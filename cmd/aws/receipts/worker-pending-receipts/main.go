package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	rf "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

/*
1) Read Receipt ID from SQS Message Attribute
 - retrieve Receipt from the database

2) Read Receipt url from SQS message
 - try to parse receipt using `receipt_fetcher` package
 - if parsed receipt data is not available, return the message to the queue and increase visibility time

3) Update Receipt
 - if parsed receipt data is available update DB Receipt with parsed data

4) Write Receipt ID to SQS queue for item categorization
*/

const RECEIPT_ITEMS_CATEGORIZE_QUEUE = "receipt_items_categorize"

var (
	//Environment
	env = os.Getenv("ENVIRONMENT")

	// AWS SQS
	session                      *aws_session.Session
	client                       *sqs.SQS
	receiptItemsCategorizeSqsUrl *string
)

func init() {
	// Initialize database
	err := db.InitializeDB()
	if err != nil {
		os.Exit(1)
	}

	// Initialize AWS session and SQS client
	options := aws_session.Options{SharedConfigState: aws_session.SharedConfigDisable}
	session = aws_session.Must(aws_session.NewSessionWithOptions(options))
	client = sqs.New(session)

	// Initialize SQS urls
	if env == "dev" {
		localStackSqsUrl := "https://localhost.localstack.cloud:4566/000000000000/receipt_items_categorize"
		receiptItemsCategorizeSqsUrl = &localStackSqsUrl
	} else {
		if urlResult, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: aws.String(RECEIPT_ITEMS_CATEGORIZE_QUEUE),
		}); err != nil {
			panic(1)
		} else {
			receiptItemsCategorizeSqsUrl = urlResult.QueueUrl
		}
	}
}

func processMessage(ctx context.Context, message events.SQSMessage) error {
	// Receipt url
	url := message.Body
	// Receipt ID
	receiptIdPtr := message.MessageAttributes["ReceiptId"].StringValue
	receiptID := *receiptIdPtr

	receipt, err := rf.Get(url)
	if err != nil {
		// If url is invalid
		if errors.Is(err, rf.ErrInvalidReceiptUrl) {
			return errors.New("invalid Receipt url")
		}

		if errors.Is(err, rf.ErrReceiptDataNotAvailable) {
			// Note: Apparently I can change message visibility parameter
			// So that I can delay visibility of message carrying receipt url
			// that could not be parsed at the time of execution
			changeVisibilityParams := &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          aws.String(message.EventSourceARN),
				ReceiptHandle:     aws.String(message.ReceiptHandle),
				VisibilityTimeout: aws.Int64(120), // 2 minutes
			}

			if _, err = client.ChangeMessageVisibilityWithContext(ctx, changeVisibilityParams); err != nil {
				return fmt.Errorf("unable to change message visibility: %v", err)
			}
		}
	}

	// Retrieve Receipt from the database
	var dbReceipt models.Receipt
	if dbErr := db.Instance.First(&dbReceipt, receiptID).Error; dbErr != nil {
		return dbErr
	}

	// Combination of Tin and LocationId should be unique
	var dbStore models.Store
	db.Instance.FirstOrCreate(&dbStore, models.Store{Tin: receipt.Store.Tin, LocationId: receipt.Store.LocationId}, func(tx *gorm.DB) error {
		dbStore.Name = receipt.Store.Name
		dbStore.LocationName = receipt.Store.LocationName
		dbStore.Address = receipt.Store.Address
		dbStore.City = receipt.Store.City

		return nil
	})

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

	var meta datatypes.JSON
	data, _ := json.Marshal(receipt.MetaData)
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}

	dbReceipt.PfrNumber = &receipt.Number
	dbReceipt.Counter = &receipt.Counter
	dbReceipt.TotalPurchaseAmount = &totalPurchaseAmount
	dbReceipt.TotalTaxAmount = &totalTaxAmount
	dbReceipt.Date = &receipt.Date
	dbReceipt.QrCode = &receipt.QrCod
	dbReceipt.Meta = &meta
	dbReceipt.Store = &dbStore
	dbReceipt.ReceiptItems = dbReceiptItems

	if dbErr := db.Instance.Save(&dbReceipt).Error; dbErr != nil {
		return dbErr
	}

	// Push Receipt ID to SQS for further processing (items categorization)
	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageBody:  aws.String(string(receiptID)),
		QueueUrl:     receiptItemsCategorizeSqsUrl,
	}

	if _, err = client.SendMessage(sqsMessageInput); err != nil {
		return err
	}

	return nil
}

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		if err := processMessage(ctx, message); err != nil {
			log.Printf("Message with ID: %s failed with error: %v\n", message.MessageId, err)
			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}