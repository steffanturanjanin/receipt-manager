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
	ReceiptRepository repositories.ReceiptRepositoryInterface
}

func NewReceiptService(receiptRepository repositories.ReceiptRepositoryInterface) *ReceiptService {
	return &ReceiptService{
		ReceiptRepository: receiptRepository,
	}
}

func (service *ReceiptService) CreateFromUrl(url string) (*dto.Receipt, error) {
	receiptData, err := receipt_fetcher.Get(url)
	if err != nil {
		return nil, errors.NewErrBadRequest(err, "Invalid receipt url.")
	}

	receiptModel, err := service.ReceiptRepository.Create(dto.ReceiptData(*receiptData))
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

	if err := service.ReceiptRepository.Create2(&r); err != nil {
		return nil, err
	}

	receiptDTO := dto.Receipt{
		ID:     r.ID,
		Status: r.Status,
	}

	return &receiptDTO, nil
}

func (service *ReceiptService) Delete(id int) error {
	return service.ReceiptRepository.Delete(id)
}

func (service *ReceiptService) GetAll(f filters.ReceiptFilters, p *pagination.Pagination) ([]dto.Receipt, error) {
	receipts := make([]dto.Receipt, 0)

	receiptModels, err := service.ReceiptRepository.GetAll(f, p)
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
			CategoryID:   item.CategoryId,
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

	return s.ReceiptRepository.Update(receipt)
}

func (s *ReceiptService) GetByPfr(pfr string) (*dto.Receipt, error) {
	receiptModel, err := s.ReceiptRepository.GetByPfr(pfr)
	if err != nil {
		return nil, err
	}

	receipt := dto.Receipt{
		ID: receiptModel.ID,
	}

	return &receipt, nil
}
