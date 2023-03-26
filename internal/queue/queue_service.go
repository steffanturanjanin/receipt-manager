package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type QueueService struct {
	sqssvc *sqs.SQS
}

func NewQueueService(sqssvc *sqs.SQS) *QueueService {
	return &QueueService{
		sqssvc: sqssvc,
	}
}

func (qs QueueService) SendMessage(qp QueueProducerInterface) error {
	sendMessageInput, err := qs.buildSendMessageInput(qp)
	if err != nil {
		return err
	}

	if _, err := qs.sqssvc.SendMessage(sendMessageInput); err != nil {
		return err
	}

	return nil
}

func (qs QueueService) buildSendMessageInput(qp QueueProducerInterface) (*sqs.SendMessageInput, error) {
	var queueName string = qp.GetQueueName()

	urlResult, err := qs.sqssvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		return nil, err
	}

	var message string = qp.GetMessage()

	sendMessageInput := &sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(0),
		MessageAttributes: qp.GetMessageAttributes(),
		MessageBody:       &message,
		QueueUrl:          urlResult.QueueUrl,
	}

	return sendMessageInput, nil
}

func (qs QueueService) SpawnWorker(w QueueWorkerInterface) {
	w.RegisterHandlers()

	go w.Consume()
}
