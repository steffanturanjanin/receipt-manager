package services

import (
	"encoding/json"
	"math"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

type ReceiptService struct {
	ReceiptRepository repositories.ReceiptRepositoryInterface
}

func NewReceiptService(receiptRepository repositories.ReceiptRepositoryInterface) *ReceiptService {
	return &ReceiptService{
		ReceiptRepository: receiptRepository,
	}
}

func (service *ReceiptService) CreateFromUrl(url string) (*dto.Receipt, error) {
	receipt, err := receipt_fetcher.Get(url)

	if err != nil {
		return nil, err
	}

	receiptModel, err := service.ReceiptRepository.Create(dto.ReceiptData(*receipt))
	if err != nil {
		return nil, err
	}

	receiptItems := make([]dto.ReceiptItem, 0)
	for _, receiptItem := range receiptModel.ReceiptItems {
		receiptItemDTO := dto.ReceiptItem{
			ID:           receiptItem.ID,
			Name:         receiptItem.Name,
			Quantity:     receiptItem.Quantity,
			Unit:         receiptItem.Unit,
			Tax:          *dto.TaxIdentifier(receiptItem.Tax).Tax(),
			SingleAmount: math.Round(float64(receiptItem.SingleAmount)) / 100,
			TotalAmount:  math.Round(float64(receiptItem.TotalAmount)) / 100,
		}

		receiptItems = append(receiptItems, receiptItemDTO)
	}

	taxes := make([]dto.Tax, 0)
	for _, taxItem := range receiptModel.Taxes {
		if tax := dto.TaxIdentifier(taxItem.TaxIdentifier).Tax(); tax != nil {
			taxes = append(taxes, *tax)
		}
	}

	meta := make(map[string]string)
	if err := json.Unmarshal(receiptModel.Meta, &meta); err != nil {
		return nil, err
	}

	receiptDTO := dto.Receipt{
		ID:                  receiptModel.ID,
		PfrNumber:           receiptModel.PfrNumber,
		Counter:             receiptModel.Counter,
		TotalPurchaseAmount: math.Round(float64(receiptModel.TotalPurchaseAmount)) / 100,
		TotalTaxAmount:      math.Round(float64(receiptModel.TotalTaxAmount)) / 100,
		ReceiptItems:        receiptItems,
		Taxes:               taxes,
		Date:                receiptModel.Date,
		QrCode:              receipt.QrCod,
		Meta:                meta,
		CreatedAt:           receiptModel.Date,
	}

	return &receiptDTO, nil
}