package queue

import (
	"context"
	"fmt"
	"os"

	"github.com/qhenkart/gosqs"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

const (
	RECEIPT_URL_QUEUE             = "receipt-url-queue"
	PROCESS_RECEIPT_HANDLER_ROUTE = "process_receipt"
)

type ReceiptUrlConsumer struct {
	gosqs.Consumer
}

func NewReceiptUrlConsumer(queueUrl string) (*ReceiptUrlConsumer, error) {
	config := gosqs.Config{
		Key:      os.Getenv("AWS_ACCESS_KEY"),
		Secret:   os.Getenv("AWS_SECRET_KEY"),
		Region:   os.Getenv("AWS_REGION"),
		QueueURL: queueUrl,
	}

	consumer, err := gosqs.NewConsumer(config, RECEIPT_URL_QUEUE)
	if err != nil {
		return nil, err
	}

	return &ReceiptUrlConsumer{consumer}, nil
}

func (c *ReceiptUrlConsumer) ProcessReceiptHandler(ctx context.Context, m gosqs.Message) error {
	u := struct {
		Url string `json:"url"`
	}{}

	if err := m.Decode(&u); err != nil {
		return err
	}

	receipt, err := receipt_fetcher.Get(u.Url)
	if err != nil {
		return err
	}

	fmt.Printf("RECEIPT PROCESSED: %+v\n", receipt)

	return nil
}

func (c *ReceiptUrlConsumer) RegisterHandlers() {
	c.RegisterHandler(PROCESS_RECEIPT_HANDLER_ROUTE, c.ProcessReceiptHandler)
}

func (c *ReceiptUrlConsumer) StartWorker() {
	c.RegisterHandlers()

	go c.Consume()
}
