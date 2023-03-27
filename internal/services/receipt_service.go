package services

import (
	"encoding/json"
	"math"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/pagination"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
)

type ReceiptService struct {
	receiptRepository repositories.ReceiptRepositoryInterface
}

func NewReceiptService(receiptRepository repositories.ReceiptRepositoryInterface) *ReceiptService {
	return &ReceiptService{
		receiptRepository: receiptRepository,
	}
}

func (service *ReceiptService) CreateFromUrl(url string) (*dto.Receipt, error) {
	receiptData, err := receipt_fetcher.Get(url)
	if err != nil {
		return nil, errors.NewErrBadRequest(err, "Invalid receipt url.")
	}

	receiptModel, err := service.receiptRepository.Create(dto.ReceiptData(*receiptData))
	if err != nil {
		return nil, err
	}

	receipt, err := receiptModel.NewReceiptDTO()
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (service *ReceiptService) CreatePendingReceipt() (*dto.Receipt, error) {
	r := models.Receipt{
		Status: models.RECEIPT_STATUS_PENDING,
	}

	if err := service.receiptRepository.Create2(&r); err != nil {
		return nil, err
	}

	receiptDTO := dto.Receipt{
		ID:     r.ID,
		Status: r.Status,
	}

	return &receiptDTO, nil
}

func (service *ReceiptService) Delete(id int) error {
	return service.receiptRepository.Delete(id)
}

func (service *ReceiptService) GetAll(f filters.ReceiptFilters, p *pagination.Pagination) ([]dto.Receipt, error) {
	receipts := make([]dto.Receipt, 0)

	receiptModels, err := service.receiptRepository.GetAll(f, p)
	if err != nil {
		return nil, err
	}

	for _, receiptModel := range receiptModels {
		receipt, err := receiptModel.NewReceiptDTO()
		if err != nil {
			return nil, err
		}

		receipts = append(receipts, *receipt)
	}

	return receipts, nil
}

func (s *ReceiptService) UpdateProcessedReceipt(r dto.ReceiptParams) error {
	receiptItems := make([]models.ReceiptItem, 0)
	for _, item := range r.ReceiptItems {
		taxId := dto.TaxIdentifierMapper[item.Tax.Identifier]

		receiptItems = append(receiptItems, models.ReceiptItem{
			Name:         item.Name,
			Unit:         item.Unit,
			Quantity:     item.Quantity,
			SingleAmount: int(math.Round(item.SingleAmount * 100)),
			TotalAmount:  int(math.Round(item.TotalAmount * 100)),
			Tax:          int(taxId),
			CategoryID:   &item.Category.Id,
		})
	}

	taxes := make([]models.Tax, 0)
	for _, tax := range r.Taxes {
		taxId := dto.TaxIdentifierMapper[tax.Identifier]
		taxes = append(taxes, models.Tax{
			TaxIdentifier: int(taxId),
		})
	}

	metaData, _ := json.Marshal(r.Meta)
	receipt := &models.Receipt{
		ID:                  *r.Id,
		Status:              models.RECEIPT_STATUS_PROCESSED,
		PfrNumber:           r.PfrNumber,
		Counter:             r.Counter,
		TotalPurchaseAmount: *r.TotalPurchaseAmount,
		TotalTaxAmount:      *r.TotalTaxAmount,
		Date:                *r.Date,
		QrCode:              r.QrCode,
		Meta:                metaData,
		Store: models.Store{
			Tin:          r.Store.Tin,
			Name:         r.Store.Name,
			LocationName: r.Store.LocationName,
			LocationId:   r.Store.LocationId,
			Address:      r.Store.Address,
			City:         r.Store.City,
		},
		ReceiptItems: receiptItems,
		Taxes:        taxes,
	}

	return s.receiptRepository.Update(receipt)
}

func (s *ReceiptService) GetByPfr(pfr string) (*dto.Receipt, error) {
	receiptModel, err := s.receiptRepository.GetByPfr(pfr)
	if err != nil {
		return nil, err
	}

	receipt := dto.Receipt{
		ID: receiptModel.ID,
	}

	return &receipt, nil
}

func (s *ReceiptService) GetById(id int) (*dto.Receipt, error) {
	receipt, err := s.receiptRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	receiptItems := make([]dto.ReceiptItem, 0)
	for _, receiptItem := range receipt.ReceiptItems {
		receiptItems = append(receiptItems, dto.ReceiptItem{
			ID:   receiptItem.ID,
			Name: receiptItem.Name,
			Category: &dto.Category{
				Id:   receiptItem.Category.ID,
				Name: receiptItem.Category.Name,
			},
			Quantity:     receiptItem.Quantity,
			Unit:         receiptItem.Unit,
			SingleAmount: math.Round(float64(receiptItem.SingleAmount)) / 100,
			TotalAmount:  math.Round(float64(receiptItem.TotalAmount)) / 100,
			Tax:          *dto.TaxIdentifier(receiptItem.Tax).Tax(),
		})
	}

	taxes := make([]dto.Tax, 0)
	for _, taxModel := range receipt.Taxes {
		tax := dto.TaxIdentifier(taxModel.TaxIdentifier).Tax()
		taxes = append(taxes, *tax)
	}

	var meta map[string]string
	json.Unmarshal(receipt.Meta, &meta)
	receiptDto := dto.Receipt{
		ID:                  receipt.ID,
		Status:              receipt.Status,
		PfrNumber:           *receipt.PfrNumber,
		Counter:             *receipt.Counter,
		TotalPurchaseAmount: math.Round(float64(receipt.TotalPurchaseAmount)) / 100,
		TotalTaxAmount:      math.Round(float64(receipt.TotalTaxAmount)) / 100,
		Date:                receipt.Date,
		QrCode:              *receipt.QrCode,
		Meta:                meta,
		ReceiptItems:        receiptItems,
		Taxes:               taxes,
		CreatedAt:           receipt.CreatedAt,
		Store: dto.Store{
			Tin:          receipt.Store.Tin,
			Name:         receipt.Store.Name,
			LocationName: receipt.Store.LocationName,
			LocationId:   receipt.Store.LocationId,
			Address:      receipt.Store.Address,
			City:         receipt.Store.City,
		},
	}

	return &receiptDto, nil
}
