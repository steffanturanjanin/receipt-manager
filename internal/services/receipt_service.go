package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/filters"
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
