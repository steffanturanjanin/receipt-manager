package repositories

import (
	"encoding/json"

	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReceiptRepositoryInterface interface {
	Create(receiptDTO dto.ReceiptData) (*models.Receipt, error)
}

type ReceiptRepository struct {
	db *gorm.DB
}

func NewReceiptRepository(db *gorm.DB) *ReceiptRepository {
	return &ReceiptRepository{
		db: db,
	}
}

func (repository *ReceiptRepository) Create(receiptDTO dto.ReceiptData) (*models.Receipt, error) {
	receiptItems := make([]models.ReceiptItem, 0)

	for _, receiptItemDTO := range receiptDTO.Items {
		receiptItem := models.ReceiptItem{
			Name:         receiptItemDTO.Name,
			Unit:         receiptItemDTO.Unit,
			Quantity:     receiptItemDTO.Quantity,
			SingleAmount: receiptItemDTO.SingleAmount.GetParas(),
			TotalAmount:  receiptItemDTO.TotalAmount.GetParas(),
			Tax:          int(dto.TaxIdentifierMapper[receiptItemDTO.Tax.Identifier]),
		}

		receiptItems = append(receiptItems, receiptItem)
	}

	taxes := make([]models.Tax, 0)
	for _, taxItem := range receiptDTO.Taxes {
		taxModel := models.Tax{
			TaxIdentifier: int(dto.TaxIdentifierMapper[taxItem.Tax.Identifier]),
		}

		taxes = append(taxes, taxModel)
	}

	metaJson, _ := json.Marshal(receiptDTO.MetaData)

	receipt := models.Receipt{
		PfrNumber:           receiptDTO.Number,
		Counter:             receiptDTO.Counter,
		TotalPurchaseAmount: receiptDTO.TotalPurchaseAmount.GetParas(),
		TotalTaxAmount:      receiptDTO.TotalTaxAmount.GetParas(),
		Date:                receiptDTO.Date,
		QrCode:              receiptDTO.QrCod,
		ReceiptItems:        receiptItems,
		Taxes:               taxes,
		Meta:                datatypes.JSON(metaJson),
	}

	if result := repository.db.Create(&receipt); result.Error != nil {
		return nil, result.Error
	}

	return &receipt, nil
}
