package controllers

import (
	"encoding/json"
	native_erors "errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"github.com/steffanturanjanin/receipt-manager/internal/pagination"
	"github.com/steffanturanjanin/receipt-manager/internal/queue"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
	"github.com/steffanturanjanin/receipt-manager/internal/validator"
)

type ReceiptController struct {
	receiptService *services.ReceiptService
	queueService   *queue.QueueService
	validator      *validator.Validator
}

func NewReceiptController(rs *services.ReceiptService, qs *queue.QueueService, v *validator.Validator) *ReceiptController {
	return &ReceiptController{
		receiptService: rs,
		queueService:   qs,
		validator:      v,
	}
}

func (controller *ReceiptController) CreateFromUrl(w http.ResponseWriter, r *http.Request) {
	url := struct{ Url string }{}

	if err := ParseBody(&url, r); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	receipt, err := controller.receiptService.CreateFromUrl(url.Url)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, receipt, http.StatusCreated)
}

func (controller *ReceiptController) CreateFromUrl2(w http.ResponseWriter, r *http.Request) {
	url := struct {
		Url string `validate:"receiptUrl" json:"url"`
	}{}

	if err := ParseBody(&url, r); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	if err := ValidateRequest(&url, controller.validator); err != nil {
		JsonErrorResponse(w, err)
		return
	}

	receiptDTO, err := controller.receiptService.CreatePendingReceipt()
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	urlWithReceiptId := struct {
		ID  uint   `json:"id"`
		Url string `json:"url"`
	}{
		ID:  receiptDTO.ID,
		Url: url.Url,
	}

	message, _ := json.Marshal(&urlWithReceiptId)

	qp := queue.NewReceiptUrlQueueProducer(string(message))
	if err := controller.queueService.SendMessage(&qp); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
	}

	JsonInfoResponse(w, "Receipt created and is set to be processed.", http.StatusOK)
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

	if err := controller.receiptService.Delete(id); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, nil, http.StatusNoContent)
}

func (controller *ReceiptController) List(w http.ResponseWriter, r *http.Request) {
	filters := filters.ReceiptFilters{}
	filters.BuildFromRequest(r)
	pagination := pagination.GetPaginationFromRequest(r)

	receipts, err := controller.receiptService.GetAll(filters, &pagination)
	if err != nil {
		JsonErrorResponse(w, err)
		return
	}

	JsonPaginatedResponse(w, receipts, pagination, http.StatusOK)
}
