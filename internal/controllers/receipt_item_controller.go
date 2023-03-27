package controllers

import (
	native_errors "errors"
	"net/http"
	"strconv"

	"github.com/steffanturanjanin/receipt-manager/internal/errors"

	"github.com/gorilla/mux"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
)

type ReceiptItemController struct {
	receiptItemService *services.ReceiptItemService
}

func NewReceiptItemController(s *services.ReceiptItemService) *ReceiptItemController {
	return &ReceiptItemController{
		receiptItemService: s,
	}
}

func (c *ReceiptItemController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		JsonErrorResponse(w, errors.NewHttpError(
			errors.NewErrBadRequest(
				native_errors.New("missing id param from url"),
				"Missing id param from url.",
			)),
		)

		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(
			errors.NewErrBadRequest(
				native_errors.New("invalid data type for id param"),
				"Invalid data type for id param.",
			)),
		)

		return
	}

	var data services.UpdateReceiptItemCategory
	if err := ParseBody(&data, r); err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	data.ReceiptItemId = id

	receiptItem, err := c.receiptItemService.UpdateCategory(data)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, receiptItem, http.StatusOK)
}
