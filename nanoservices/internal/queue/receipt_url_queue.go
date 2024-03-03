package queue

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/qhenkart/gosqs"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/pkg/receipt-fetcher"
)

const (
	QUEUE_RECEIPT_URL                             = "receipt-url-queue"
	HANDLER_QUEUE_RECEIPT_URL_PROCESS_RECEIPT_URL = "process_receipt_url"
)

// Queue Producer
type ReceiptUrlQueueProducer struct {
	message string
}

func NewReceiptUrlQueueProducer(message string) ReceiptUrlQueueProducer {
	return ReceiptUrlQueueProducer{
		message: message,
	}
}

func (qp *ReceiptUrlQueueProducer) GetMessage() string {
	return qp.message
}

func (qp *ReceiptUrlQueueProducer) GetQueueName() string {
	return QUEUE_RECEIPT_URL
}

func (qp *ReceiptUrlQueueProducer) GetMessageAttributes() map[string]*sqs.MessageAttributeValue {
	return map[string]*sqs.MessageAttributeValue{
		"route": {
			DataType:    aws.String("String"),
			StringValue: aws.String(HANDLER_QUEUE_RECEIPT_URL_PROCESS_RECEIPT_URL),
		},
	}
}

// Queue Worker
type ReceiptUrlQueueWorker struct {
	gosqs.Consumer
	queueService *QueueService
}

func NewReceiptUrlQueueWorker(qs *QueueService) (*ReceiptUrlQueueWorker, error) {
	var queueName string = QUEUE_RECEIPT_URL
	queueUrl, err := qs.sqssvc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		return nil, err
	}

	config := gosqs.Config{
		Key:      os.Getenv("AWS_ACCESS_KEY"),
		Secret:   os.Getenv("AWS_SECRET_KEY"),
		Region:   os.Getenv("AWS_REGION"),
		QueueURL: *queueUrl.QueueUrl,
	}

	consumer, err := gosqs.NewConsumer(config, QUEUE_RECEIPT_URL)
	if err != nil {
		return nil, err
	}

	return &ReceiptUrlQueueWorker{consumer, qs}, nil
}

type UrlWithReceiptId struct {
	ID  uint   `json:"id"`
	Url string `json:"url"`
}

func (w *ReceiptUrlQueueWorker) ProcessReceiptUrlHandler(ctx context.Context, m gosqs.Message) error {
	var urlWithReceiptId UrlWithReceiptId

	if err := m.Decode(&urlWithReceiptId); err != nil {
		return err
	}

	receipt, err := receipt_fetcher.Get(urlWithReceiptId.Url)
	if err != nil {
		return err
	}

	parsedReceipt := parsedReceiptQueueData{
		ID:            urlWithReceiptId.ID,
		ParsedReceipt: dto.ReceiptData(*receipt),
	}

	qp := NewReceiptParsedQueueProducer(parsedReceipt)
	if err := w.queueService.SendMessage(&qp); err != nil {
		return err
	}

	return nil
}

func (w *ReceiptUrlQueueWorker) RegisterHandlers() {
	w.RegisterHandler(HANDLER_QUEUE_RECEIPT_URL_PROCESS_RECEIPT_URL, w.ProcessReceiptUrlHandler)
}
