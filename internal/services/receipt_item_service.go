package services

import (
	native_errors "errors"
	"math"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	"github.com/steffanturanjanin/receipt-manager/internal/utils"
)

type ReceiptItemService struct {
	receiptItemRepository repositories.ReceiptItemRepositoryInterface
	categoryService       *CategoryService
}

func NewReceiptItemService(r repositories.ReceiptItemRepositoryInterface, cs *CategoryService) *ReceiptItemService {
	return &ReceiptItemService{
		receiptItemRepository: r,
		categoryService:       cs,
	}
}

type UpdateReceiptItemCategory struct {
	ReceiptItemId int `json:"receipt_item_id"`
	CategoryId    int `json:"category_id"`
}

func (s *ReceiptItemService) UpdateCategory(data UpdateReceiptItemCategory) (*dto.ReceiptItem, error) {
	ids, err := s.categoryService.GetIds()
	if err != nil || !utils.InSlice(ids, data.CategoryId) {
		return nil, errors.NewErrBadRequest(native_errors.New("invalid category id"), "Invalid category id.")
	}

	receiptItem, err := s.receiptItemRepository.FindById(data.ReceiptItemId)
	if err != nil {
		return nil, err
	}

	categoryId := new(uint)
	*categoryId = uint(data.CategoryId)

	receiptItem.CategoryID = categoryId

	if err := s.receiptItemRepository.Update(receiptItem); err != nil {
		return nil, err
	}

	var category *dto.Category
	if receiptItem.Category != nil {
		category.Id = receiptItem.Category.ID
		category.Name = receiptItem.Category.Name
	}

	receiptItemDto := dto.ReceiptItem{
		ID:       receiptItem.ID,
		Name:     receiptItem.Name,
		Category: category,
		Unit:     receiptItem.Unit,
		//Tax:          *dto.TaxIdentifier(receiptItem.Tax).Tax(),
		SingleAmount: math.Round(float64(receiptItem.SingleAmount)) / 100,
		TotalAmount:  math.Round(float64(receiptItem.TotalAmount)) / 100,
	}

	return &receiptItemDto, nil
}
