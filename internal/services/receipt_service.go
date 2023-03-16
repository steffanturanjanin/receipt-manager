package services

import (
	"encoding/json"
	"math"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/pagination"
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
		return nil, errors.NewErrBadRequest(err, "Invalid receipt url.")
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
		ID: receiptModel.ID,
		Store: dto.Store{
			Tin:          receiptModel.Store.Tin,
			Name:         receiptModel.Store.Name,
			LocationName: receiptModel.Store.LocationName,
			LocationId:   receiptModel.Store.LocationId,
			Address:      receiptModel.Store.Address,
			City:         receiptModel.Store.City,
		},
		PfrNumber:           receiptModel.PfrNumber,
		Counter:             receiptModel.Counter,
		TotalPurchaseAmount: math.Round(float64(receiptModel.TotalPurchaseAmount)) / 100,
		TotalTaxAmount:      math.Round(float64(receiptModel.TotalTaxAmount)) / 100,
		ReceiptItems:        receiptItems,
		Taxes:               taxes,
		Date:                receiptModel.Date,
		QrCode:              receiptModel.QrCode,
		Meta:                meta,
		CreatedAt:           receiptModel.Date,
	}

	return &receiptDTO, nil
}

func (service *ReceiptService) Delete(id int) error {
	return service.ReceiptRepository.Delete(id)
}

func (service *ReceiptService) GetAll(p *pagination.Pagination) ([]dto.Receipt, error) {
	receipts := make([]dto.Receipt, 0)

	receiptModels, err := service.ReceiptRepository.GetAll(p)
	if err != nil {
		return nil, err
	}

	for _, receiptModel := range receiptModels {
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
			ID: receiptModel.ID,
			Store: dto.Store{
				Tin:          receiptModel.Store.Tin,
				Name:         receiptModel.Store.Name,
				LocationName: receiptModel.Store.LocationName,
				LocationId:   receiptModel.Store.LocationId,
				Address:      receiptModel.Store.Address,
				City:         receiptModel.Store.City,
			},
			PfrNumber:           receiptModel.PfrNumber,
			Counter:             receiptModel.Counter,
			TotalPurchaseAmount: math.Round(float64(receiptModel.TotalPurchaseAmount)) / 100,
			TotalTaxAmount:      math.Round(float64(receiptModel.TotalTaxAmount)) / 100,
			ReceiptItems:        receiptItems,
			Taxes:               taxes,
			Date:                receiptModel.Date,
			QrCode:              receiptModel.QrCode,
			Meta:                meta,
			CreatedAt:           receiptModel.Date,
		}

		receipts = append(receipts, receiptDTO)
	}

	return receipts, nil
}
