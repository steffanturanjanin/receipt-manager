package queue

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type ReceiptUrlProducer struct {
	sqs      *sqs.SQS
	queueUrl string
}

func NewReceiptUrlProducer(sqs *sqs.SQS, queueUrl string) ReceiptUrlProducer {
	return ReceiptUrlProducer{
		sqs:      sqs,
		queueUrl: queueUrl,
	}
}

func (p ReceiptUrlProducer) Produce(url string) error {
	u := struct {
		Url string `json:"url"`
	}{Url: url}

	msg, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = p.sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"route": {
				DataType:    aws.String("String"),
				StringValue: aws.String(PROCESS_RECEIPT_HANDLER_ROUTE),
			},
		},
		MessageBody: aws.String(string(msg)),
		QueueUrl:    &p.queueUrl,
	})

	if err != nil {
		return err
	}

	return nil
}
