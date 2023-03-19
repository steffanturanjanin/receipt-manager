package repositories

import (
	"encoding/json"
	native_errors "errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/pagination"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReceiptRepositoryInterface interface {
	GetAll(f filters.ReceiptFilters, p *pagination.Pagination) ([]models.Receipt, error)
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

func (repository *ReceiptRepository) Create(receiptData dto.ReceiptData) (*models.Receipt, error) {
	var store models.Store
	repository.db.FirstOrCreate(&store, models.Store{
		Tin:          receiptData.Store.Tin,
		Name:         receiptData.Store.Name,
		LocationId:   receiptData.Store.LocationId,
		LocationName: receiptData.Store.Name,
		Address:      receiptData.Store.Address,
		City:         receiptData.Store.City,
	})

	receiptItems := make([]models.ReceiptItem, 0)

	for _, receiptItemDTO := range receiptData.Items {
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
	for _, taxItem := range receiptData.Taxes {
		taxModel := models.Tax{
			TaxIdentifier: int(dto.TaxIdentifierMapper[taxItem.Tax.Identifier]),
		}

		taxes = append(taxes, taxModel)
	}

	metaJson, _ := json.Marshal(receiptData.MetaData)

	receipt := models.Receipt{
		StoreID:             store.Tin,
		PfrNumber:           receiptData.Number,
		Counter:             receiptData.Counter,
		TotalPurchaseAmount: receiptData.TotalPurchaseAmount.GetParas(),
		TotalTaxAmount:      receiptData.TotalTaxAmount.GetParas(),
		Date:                receiptData.Date,
		QrCode:              receiptData.QrCod,
		ReceiptItems:        receiptItems,
		Taxes:               taxes,
		Meta:                datatypes.JSON(metaJson),
	}

	if result := repository.db.Create(&receipt).Preload("Store"); result.Error != nil {
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

func (repository *ReceiptRepository) GetAll(f filters.ReceiptFilters, p *pagination.Pagination) ([]models.Receipt, error) {
	var receipts []models.Receipt

	baseQuery := repository.db.Model(models.Receipt{})
	fmt.Printf("FILTERS_LIST %+v\n", f.FiltersList)

	filteredQuery := f.ApplyFilters(baseQuery)
	paginatedQuery := Paginate(filteredQuery, p, models.Receipt{})

	result := paginatedQuery.Preload("Store").Preload("ReceiptItems").Find(&receipts)

	if result.Error != nil {
		return nil, result.Error
	}

	return receipts, nil
}
