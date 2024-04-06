package transport

import (
	"fmt"

	"github.com/steffanturanjanin/receipt-manager/internal/models"
)

type ReceiptItemResponse struct {
	ID           int               `json:"id"`
	ReceiptId    int               `json:"receiptId"`
	Name         string            `json:"name"`
	Unit         string            `json:"unit"`
	Quantity     float64           `json:"quantity"`
	SingleAmount string            `json:"singleAmount"`
	TotalAmount  string            `json:"totalAmount"`
	Category     *CategoryResponse `json:"category"`
	Tax          *TaxResponse      `json:"tax"`
}

type ReceiptItemTransformer struct{}

func (t ReceiptItemTransformer) TransformSingle(model models.ReceiptItem) ReceiptItemResponse {
	receiptItem := ReceiptItemResponse{}

	var category *CategoryResponse
	if model.Category != nil {
		transformer := CategoryTransformer{}
		category = new(CategoryResponse)
		*category = transformer.TransformSingle(*model.Category)
	}
	var tax *TaxResponse
	if model.Tax != nil {
		transformer := TaxTransformer{}
		tax = new(TaxResponse)
		*tax = transformer.TransformSingle(*model.Tax)
	}

	receiptItem.ID = int(model.ID)
	receiptItem.ReceiptId = int(model.ReceiptID)
	receiptItem.Name = model.Name
	receiptItem.Unit = model.Unit
	receiptItem.Quantity = model.Quantity
	receiptItem.SingleAmount = fmt.Sprintf("%.2f", float64(model.SingleAmount)/100)
	receiptItem.TotalAmount = fmt.Sprintf("%.2f", float64(model.TotalAmount)/100)
	receiptItem.Category = category
	receiptItem.Tax = tax

	return receiptItem
}
