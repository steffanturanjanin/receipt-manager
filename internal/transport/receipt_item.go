package transport

type ReceiptItemResponse struct {
	ID           int               `json:"id"`
	ReceiptId    int               `json:"receiptId"`
	Name         string            `json:"name"`
	Unit         string            `json:"unit"`
	Quantity     float64           `json:"quantity"`
	SingleAmount int               `json:"singleAmount"`
	TotalAmount  int               `json:"totalAmount"`
	Category     *CategoryResponse `json:"category"`
	Tax          *TaxResponse      `json:"tax"`
}
