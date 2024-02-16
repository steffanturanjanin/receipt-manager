package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

func init() {
	// Initialize database
	if err := db.InitializeDB(); err != nil {
		os.Exit(1)
	}
}

func processMessage(ctx context.Context, message events.SQSMessage) error {
	receiptIdIntValue, err := strconv.Atoi(message.Body)
	if err != nil {
		return err
	}

	receiptId := uint(receiptIdIntValue)

	var receipt models.Receipt
	if dbErr := db.Instance.Preload("ReceiptItems").First(&receipt, receiptId).Error; dbErr != nil {
		return dbErr
	}

	//Categorize receipt items...

	log.Printf("Receipt: %v\n", receipt)

	//Update receipt items categories...

	return nil
}

func handler(ctx context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		err := processMessage(ctx, message)
		if err != nil {
			log.Printf("Error while trying to categorize receipt items: %+v\n", err)
			continue
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
