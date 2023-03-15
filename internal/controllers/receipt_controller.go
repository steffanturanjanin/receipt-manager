package controllers

import (
	native_erors "errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
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
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	receipt, err := controller.ReceiptService.CreateFromUrl(url.Url)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, receipt, http.StatusCreated)
}

func (controller *ReceiptController) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := params["id"]; !ok {
		JsonErrorResponse(w, errors.NewHttpError(native_erors.New("missing parameter id from the url")))
		return
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	if err := controller.ReceiptService.Delete(id); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, nil, http.StatusNoContent)
}
