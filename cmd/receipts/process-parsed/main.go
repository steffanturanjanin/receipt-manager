package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"gorm.io/datatypes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

const (
	RECEIPT_PARSED_QUEUE = "receipt_parsed"
	RECEIPT_ITEMS_QUEUE  = "receipt_items"
)

var (
	SessionOptions *session.Options
	Session        *session.Session
	SqsService     *sqs.SQS
)

func init() {
	err := database.InitializeDB()
	if err != nil {
		panic(1)
	}

	SessionOptions := session.Options{SharedConfigState: session.SharedConfigDisable}
	Session = session.Must(session.NewSessionWithOptions(SessionOptions))
	SqsService = sqs.New(Session)
}

func processMessage(ctx context.Context, message events.SQSMessage) error {
	userId, err := strconv.Atoi(*message.MessageAttributes["UserId"].StringValue)
	if err != nil {
		return err
	}

	parsedReceipt := &receipt_fetcher.Receipt{}
	err = json.Unmarshal([]byte(message.Body), parsedReceipt)
	log.Printf("PARSED RECEIPT %+v\n", parsedReceipt)
	if err != nil {
		return err
	}

	store := models.Store{
		Tin:          parsedReceipt.Store.Tin,
		Name:         parsedReceipt.Store.Name,
		LocationId:   parsedReceipt.Store.LocationId,
		LocationName: parsedReceipt.Store.LocationName,
		Address:      parsedReceipt.Store.Address,
		City:         parsedReceipt.Store.City,
	}

	//Save store
	// dbResult := database.Instance.Create(&store)
	// if dbResult.Error != nil {
	// 	return dbResult.Error
	// }

	receiptItems := make([]models.ReceiptItem, 0)
	for _, receiptItemDto := range parsedReceipt.Items {
		receiptItems = append(receiptItems, models.ReceiptItem{
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

	metaData, _ := json.Marshal(parsedReceipt.MetaData)
	receipt := &models.Receipt{
		UserID:              uint(userId),
		Status:              models.RECEIPT_STATUS_PENDING,
		PfrNumber:           parsedReceipt.Number,
		Counter:             parsedReceipt.Counter,
		TotalPurchaseAmount: parsedReceipt.TotalPurchaseAmount.GetParas(),
		TotalTaxAmount:      parsedReceipt.TotalTaxAmount.GetParas(),
		Date:                parsedReceipt.Date,
		QrCode:              parsedReceipt.QrCod,
		Meta:                datatypes.JSON(metaData),
		Store:               store,
		ReceiptItems:        receiptItems,
	}

	// Write receipt to database
	dbResult := database.Instance.Create(receipt)
	if dbResult.Error != nil {
		return dbResult.Error
	}

	// Send message to `receipt_items` SQS queue to categorize receipt items
	urlResult, err := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_ITEMS_QUEUE),
	})
	if err != nil {
		return err
	}

	// Serialize receipt items to json string
	serializedReceiptItems, err := json.Marshal(receipt.ReceiptItems)
	if err != nil {
		return err
	}

	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ReceiptId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint(receipt.ID)),
			},
		},
		MessageBody: aws.String(string(serializedReceiptItems)),
		QueueUrl:    urlResult.QueueUrl,
	}

	_, err = SqsService.SendMessage(sqsMessageInput)
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
