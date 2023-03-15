package controllers

import (
	"net/http"

	"github.com/steffanturanjanin/receipt-manager/internal/services"
)

type ReceiptController struct {
	ReceiptService *services.ReceiptService
}

func NewReceiptController(receiptService *services.ReceiptService) *ReceiptController {
	return &ReceiptController{
		ReceiptService: receiptService,
	}
}

func (controller *ReceiptController) CreateFromUrl(w http.ResponseWriter, r *http.Request) {
	url := struct{ Url string }{}

	if err := ParseBody(&url, r); err != nil {
		JsonErrorResponse(w, err)
		return
	}

	receipt, err := controller.ReceiptService.CreateFromUrl(url.Url)
	if err != nil {
		JsonErrorResponse(w, err)
		return
	}

	JsonResponse(w, receipt, http.StatusCreated)
}
