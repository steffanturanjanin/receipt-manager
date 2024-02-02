package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

const (
	RECEIPT_PARSED_QUEUE = "receipt_parsed"
)

var (
	SessionOptions *session.Options
	Session        *session.Session
	SqsService     *sqs.SQS
)

func init() {
	SessionOptions := session.Options{SharedConfigState: session.SharedConfigDisable}
	Session = session.Must(session.NewSessionWithOptions(SessionOptions))
	SqsService = sqs.New(Session)
}

func processMessage(ctx context.Context, message events.SQSMessage) error {
	userId := message.Attributes["UserId"]
	receiptUrl := message.Body

	receipt, err := receipt_fetcher.Get(receiptUrl)
	if err != nil {
		// If url is invalid
		if errors.Is(err, receipt_fetcher.ErrInvalidReceiptUrl) {
			return errors.New("invalid Receipt url")
		}

		if errors.Is(err, receipt_fetcher.ErrReceiptDataNotAvailable) {
			// Note: Apparently I can change message visibility parameter
			// So that I can delay visibility of message carrying receipt url
			// that could not be parsed at the time of execution
			changeVisibilityParams := &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          aws.String(message.EventSourceARN),
				ReceiptHandle:     aws.String(message.ReceiptHandle),
				VisibilityTimeout: aws.Int64(120), // 2 minutes
			}

			_, err = SqsService.ChangeMessageVisibilityWithContext(ctx, changeVisibilityParams)
			if err != nil {
				return fmt.Errorf("unable to change message visibility: %v", err)
			}
		}
	}

	receiptParsed, err := json.Marshal(receipt)
	if err != nil {
		return err
	}

	urlResult, err := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_PARSED_QUEUE),
	})

	if err != nil {
		return err
	}

	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint(userId)),
			},
		},
		MessageBody: aws.String(string(receiptParsed)),
		QueueUrl:    urlResult.QueueUrl,
	}

	_, err = SqsService.SendMessage(sqsMessageInput)
	if err != nil {
		panic(err)
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
