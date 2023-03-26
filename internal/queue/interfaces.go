package queue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/qhenkart/gosqs"
)

type QueueProducerInterface interface {
	GetMessage() string
	GetQueueName() string
	GetMessageAttributes() map[string]*sqs.MessageAttributeValue
}

type QueueWorkerInterface interface {
	gosqs.Consumer
	RegisterHandlers()
}
