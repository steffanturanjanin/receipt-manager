package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/datatypes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	receiptFetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

const (
	RECEIPT_ITEMS_QUEUE = "receipt_items"
)

var (
	session *awsSession.Session
	client  *sqs.SQS

	receiptItemsQueueUrl *string
)

func init() {
	err := database.InitializeDB()
	if err != nil {
		panic(1)
	}

	sessionOptions := awsSession.Options{SharedConfigState: awsSession.SharedConfigDisable}
	session = awsSession.Must(awsSession.NewSessionWithOptions(sessionOptions))
	client = sqs.New(session)

	if urlResult, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_ITEMS_QUEUE),
	}); err != nil {
		panic(1)
	} else {
		receiptItemsQueueUrl = urlResult.QueueUrl
	}
}

func processMessage(ctx context.Context, message events.SQSMessage) error {
	userId, err := strconv.Atoi(*message.MessageAttributes["UserId"].StringValue)
	if err != nil {
		return err
	}

	receipt := &receiptFetcher.Receipt{}
	err = json.Unmarshal([]byte(message.Body), receipt)
	if err != nil {
		return err
	}

	// Check if user has scanned this receipt before
	var dbReceipt *models.Receipt
	userID := uint(userId)

	database.Instance.Where(&models.Receipt{UserID: &userID, PfrNumber: &receipt.Number}).First(&dbReceipt)
	if dbReceipt != nil {
		return errors.New("user has already scanned this receipt")
	}

	// TODO: combination oo Tin and LocationId should be unique
	dbStore := models.Store{
		Tin:          receipt.Store.Tin,
		Name:         receipt.Store.Name,
		LocationId:   receipt.Store.LocationId,
		LocationName: receipt.Store.LocationName,
		Address:      receipt.Store.Address,
		City:         receipt.Store.City,
	}

	dbReceiptItems := make([]models.ReceiptItem, 0)
	for _, receiptItemDto := range receipt.Items {
		dbReceiptItems = append(dbReceiptItems, models.ReceiptItem{
			Name:         receiptItemDto.Name,
			Unit:         receiptItemDto.Unit,
			Quantity:     receiptItemDto.Quantity,
			SingleAmount: receiptItemDto.SingleAmount.GetParas(),
			TotalAmount:  receiptItemDto.TotalAmount.GetParas(),
			Tax: &models.Tax{
				Identifier: receiptItemDto.Tax.Identifier,
				Name:       receiptItemDto.Tax.Name,
				Rate:       receiptItemDto.Tax.Rate,
			},
		})
	}

	totalPurchaseAmount := receipt.TotalPurchaseAmount.GetParas()
	totalTaxAmount := receipt.TotalTaxAmount.GetParas()

	metaData, _ := json.Marshal(receipt.MetaData)
	var metaJson datatypes.JSON
	json.Unmarshal(metaData, &metaJson)

	dbReceipt = &models.Receipt{
		UserID:              &userID,
		Status:              models.RECEIPT_STATUS_PENDING,
		PfrNumber:           &receipt.Number,
		Counter:             &receipt.Counter,
		TotalPurchaseAmount: &totalPurchaseAmount,
		TotalTaxAmount:      &totalTaxAmount,
		Date:                &receipt.Date,
		QrCode:              &receipt.QrCod,
		Meta:                &metaJson,
		Store:               &dbStore,
		ReceiptItems:        dbReceiptItems,
	}

	// Write receipt to database
	dbResult := database.Instance.Create(&dbReceipt)
	if dbResult.Error != nil {
		return dbResult.Error
	}

	// Serialize receipt items to json string
	serializedReceiptItems, err := json.Marshal(dbReceipt.ReceiptItems)
	if err != nil {
		return err
	}

	// Send message to `receipt_items` SQS queue to categorize receipt items
	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ReceiptId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint(dbReceipt.ID)),
			},
		},
		MessageBody: aws.String(string(serializedReceiptItems)),
		QueueUrl:    receiptItemsQueueUrl,
	}

	_, err = client.SendMessage(sqsMessageInput)
	if err != nil {
		return err
	}

	return nil
}

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		err := processMessage(ctx, message)
		if err != nil {
			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
