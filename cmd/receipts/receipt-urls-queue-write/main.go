package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type UnprocessedReceipt struct {
	UserId     string `json:"userId"`
	ReceiptUrl string `json:"receiptUrl"`
}

const (
	RECEIPT_URLS_QUEUE = "receipt_urls"
)

var (
	Session    *session.Session
	SqsService *sqs.SQS
)

func init() {
	sessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	Session = session.Must(session.NewSessionWithOptions(sessionOptions))
	SqsService = sqs.New(Session)
}

func Handler(ctx context.Context, event UnprocessedReceipt) error {
	urlResult, err := SqsService.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RECEIPT_URLS_QUEUE),
	})

	if err != nil {
		return err
	}

	sqsMessageInput := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"UserId": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprint(event.UserId)),
			},
		},
		MessageBody: aws.String(event.ReceiptUrl),
		QueueUrl:    urlResult.QueueUrl,
	}

	_, err = SqsService.SendMessage(sqsMessageInput)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
