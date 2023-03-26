package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/qhenkart/gosqs"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
)

const (
	QUEUE_RECEIPT_PARSED              = "receipt-parsed-queue"
	HANDLER_RECEIPT_PARSED_CATEGORIZE = "receipt-parsed-categorize"
)

// Queue Producer
type ReceiptParsedQueueProducer struct {
	data parsedReceiptQueueData
}

type parsedReceiptQueueData struct {
	ID            uint            `json:"id"`
	ParsedReceipt dto.ReceiptData `json:"parsed_receipt"`
}

func NewReceiptParsedQueueProducer(data parsedReceiptQueueData) ReceiptParsedQueueProducer {
	return ReceiptParsedQueueProducer{
		data: data,
	}
}

func (qp *ReceiptParsedQueueProducer) GetMessage() string {
	message, _ := json.Marshal(&qp.data)

	return string(message)
}

func (qp *ReceiptParsedQueueProducer) GetQueueName() string {
	return QUEUE_RECEIPT_PARSED
}

func (qp *ReceiptParsedQueueProducer) GetMessageAttributes() map[string]*sqs.MessageAttributeValue {
	return map[string]*sqs.MessageAttributeValue{
		"route": {
			DataType:    aws.String("String"),
			StringValue: aws.String(HANDLER_RECEIPT_PARSED_CATEGORIZE),
		},
	}
}

// Queue Worker
type ParsedReceiptQueueWorker struct {
	gosqs.Consumer
	qs  *QueueService
	rs  *services.ReceiptService
	cs  *services.CategoryService
	czs *services.CategorizeService
}

func NewParsedReceiptQueueWorker(qs *QueueService, rs *services.ReceiptService, cs *services.CategoryService) (*ParsedReceiptQueueWorker, error) {
	var queueName string = QUEUE_RECEIPT_PARSED
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

	consumer, err := gosqs.NewConsumer(config, QUEUE_RECEIPT_PARSED)
	if err != nil {
		return nil, err
	}

	return &ParsedReceiptQueueWorker{consumer, qs, rs, cs, &services.CategorizeService{}}, nil
}

func (w *ParsedReceiptQueueWorker) CategoriezeReceiptItemsHandler(ctx context.Context, m gosqs.Message) error {
	var parsedReceiptData parsedReceiptQueueData
	if err := m.Decode(&parsedReceiptData); err != nil {
		return err
	}

	if r, _ := w.rs.GetByPfr(parsedReceiptData.ParsedReceipt.Number); r != nil {
		w.rs.Delete(int(parsedReceiptData.ID))
		return nil
	}

	receiptItemInList := make(services.ReceiptItemInList, 0)
	for _, receiptItem := range parsedReceiptData.ParsedReceipt.Items {
		receiptItemInList = append(receiptItemInList, services.ReceiptItemIn{
			Name: receiptItem.Name,
		})
	}

	categories, err := w.cs.GetAll()
	if err != nil {
		return err
	}

	categoryInList := make(services.CategoryInList, 0)
	for _, category := range categories {
		categoryInList = append(categoryInList, services.CategoryIn{
			Id:   category.Id,
			Name: category.Name,
		})
	}

	categorizedReceiptItemsMap, err := w.czs.Categorize(receiptItemInList, categoryInList)
	if err != nil {
		return err
	}

	var receiptData parsedReceiptQueueData
	if err := m.Decode(&receiptData); err != nil {
		return err
	}

	receiptItems := make([]dto.ReceiptItem, 0)
	for _, receiptItem := range receiptData.ParsedReceipt.Items {
		categoryId := categorizedReceiptItemsMap[receiptItem.Name]
		receiptItems = append(receiptItems, dto.ReceiptItem{
			Name:         receiptItem.Name,
			CategoryId:   categoryId,
			Quantity:     receiptItem.Quantity,
			Unit:         receiptItem.Unit,
			Tax:          dto.Tax(receiptItem.Tax),
			SingleAmount: receiptItem.SingleAmount.GetFloat(),
			TotalAmount:  receiptItem.TotalAmount.GetFloat(),
		})
	}

	taxItems := make([]dto.Tax, 0)
	for _, taxItem := range receiptData.ParsedReceipt.Taxes {
		taxItems = append(taxItems, dto.Tax(taxItem.Tax))
	}

	totalPurchaseAmount := receiptData.ParsedReceipt.TotalPurchaseAmount.GetParas()
	totalTaxAmount := receiptData.ParsedReceipt.TotalTaxAmount.GetParas()
	receiptParasms := dto.ReceiptParams{
		Id:                  &receiptData.ID,
		PfrNumber:           &receiptData.ParsedReceipt.Number,
		Counter:             &receiptData.ParsedReceipt.Counter,
		TotalPurchaseAmount: &totalPurchaseAmount,
		TotalTaxAmount:      &totalTaxAmount,
		Date:                &receiptData.ParsedReceipt.Date,
		QrCode:              &receiptData.ParsedReceipt.QrCod,
		Meta:                map[string]string(receiptData.ParsedReceipt.MetaData),
		Store: &dto.Store{
			Tin:          receiptData.ParsedReceipt.Store.Tin,
			Name:         receiptData.ParsedReceipt.Store.Name,
			LocationName: receiptData.ParsedReceipt.Store.LocationName,
			LocationId:   receiptData.ParsedReceipt.Store.LocationId,
			Address:      receiptData.ParsedReceipt.Store.Address,
			City:         receiptData.ParsedReceipt.Store.City,
		},
		ReceiptItems: receiptItems,
		Taxes:        taxItems,
	}

	if err := w.rs.UpdateProcessedReceipt(receiptParasms); err != nil {
		return err
	}

	return nil
}

func (w *ParsedReceiptQueueWorker) RegisterHandlers() {
	w.RegisterHandler(HANDLER_RECEIPT_PARSED_CATEGORIZE, w.CategoriezeReceiptItemsHandler)
}
