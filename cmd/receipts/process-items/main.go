package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

func processMessage(ctx context.Context, message events.SQSMessage) error {
	receiptId, err := strconv.Atoi(*message.MessageAttributes["ReceiptId"].StringValue)
	if err != nil {
		return err
	}

	var receiptItems []models.ReceiptItem
	err = json.Unmarshal([]byte(message.Body), &receiptItems)
	if err != nil {
		return err
	}

	//Categorize receipt items...

	log.Printf("Receipt id: %d\n", receiptId)

	//Update receipt items categories...

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
