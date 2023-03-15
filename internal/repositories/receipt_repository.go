package repositories

import (
	"encoding/json"
	native_errors "errors"

	"github.com/go-sql-driver/mysql"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReceiptRepositoryInterface interface {
	Create(receiptDTO dto.ReceiptData) (*models.Receipt, error)
	Delete(id int) error
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
	var store models.Store
	repository.db.FirstOrCreate(&store, models.Store{
		Tin:          receiptDTO.Store.Tin,
		Name:         receiptDTO.Store.Name,
		LocationId:   receiptDTO.Store.LocationId,
		LocationName: receiptDTO.Store.Name,
		Address:      receiptDTO.Store.Address,
		City:         receiptDTO.Store.City,
	})

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
		StoreID:             store.Tin,
		Store:               store,
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
		err := result.Error
		var mysqlErr *mysql.MySQLError
		if native_errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			err = errors.NewErrDuplicateEntry(err, "This receipt has already been processed.")
		}

		return nil, err
	}

	return &receipt, nil
}

func (repository *ReceiptRepository) Delete(id int) error {
	var receipt models.Receipt

	if result := repository.db.First(&receipt, id); result.Error != nil {
		err := result.Error
		if native_errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewErrResourceNotFound(err, "Receipt failed to be deleted because it does not exist.")
		}

		return err
	}

	if result := repository.db.Delete(&receipt); result.Error != nil {
		return result.Error
	}

	return nil
}
