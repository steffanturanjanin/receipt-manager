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
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/receipt-fetcher"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReceiptRepositoryInterface interface {
	GetAll(f filters.ReceiptFilters, p *pagination.Pagination) ([]models.Receipt, error)
	GetByPfr(string) (*models.Receipt, error)
	GetById(int) (*models.Receipt, error)
	Create(receiptDTO dto.ReceiptData) (*models.Receipt, error)
	Create2(*models.Receipt) error
	CreatePendingFromDto(receiptDto receipt_fetcher.Receipt, userId uint, storeId string) (*models.Receipt, error)
	Delete(id int) error
	Update(receipt *models.Receipt) error
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
			//Tax:          int(dto.TaxIdentifierMapper[receiptItemDTO.Tax.Identifier]),
		}

		receiptItems = append(receiptItems, receiptItem)
	}

	// taxes := make([]models.Tax, 0)
	// for _, taxItem := range receiptData.Taxes {
	// 	taxModel := models.Tax{
	// 		TaxIdentifier: int(dto.TaxIdentifierMapper[taxItem.Tax.Identifier]),
	// 	}

	// 	taxes = append(taxes, taxModel)
	// }

	metaJson, _ := json.Marshal(receiptData.MetaData)

	receipt := models.Receipt{
		PfrNumber:           receiptData.Number,
		Counter:             receiptData.Counter,
		TotalPurchaseAmount: receiptData.TotalPurchaseAmount.GetParas(),
		TotalTaxAmount:      receiptData.TotalTaxAmount.GetParas(),
		Date:                receiptData.Date,
		QrCode:              receiptData.QrCod,
		ReceiptItems:        receiptItems,
		//Taxes:               taxes,
		Meta: datatypes.JSON(metaJson),
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

func (receiptRepository *ReceiptRepository) Create2(r *models.Receipt) error {
	return receiptRepository.db.Create(&r).Error
}

func (receiptRepository *ReceiptRepository) CreatePendingFromDto(
	receiptDto receipt_fetcher.Receipt,
	userId uint,
	storeId string,
) (*models.Receipt, error) {
	var receipt *models.Receipt
	receiptRepository.db.Where(&models.Receipt{PfrNumber: receiptDto.Number, UserID: userId}).First(receipt)

	if receipt != nil {
		return nil, errors.NewErrResourceNotFound(native_errors.New("receipt already exists"), "receipt already exists")
	}

	// taxes := make([]models.Tax, 0)
	// for _, taxItem := range receiptDto.Taxes {
	// 	taxModel := models.Tax{
	// 		TaxIdentifier: int(dto.TaxIdentifierMapper[taxItem.Tax.Identifier]),
	// 	}

	// 	taxes = append(taxes, taxModel)
	// }

	receiptModelMeta, _ := json.Marshal(map[string]string(receiptDto.MetaData))
	receipt = &models.Receipt{
		StoreID:             storeId,
		Status:              models.RECEIPT_STATUS_PENDING,
		PfrNumber:           receiptDto.Number,
		Counter:             receiptDto.Counter,
		TotalPurchaseAmount: receiptDto.TotalPurchaseAmount.GetParas(),
		TotalTaxAmount:      receiptDto.TotalTaxAmount.GetParas(),
		Date:                receiptDto.Date,
		QrCode:              receiptDto.QrCod,
		Meta:                receiptModelMeta,
		//Taxes:               taxes,
	}

	receiptRepository.db.Create(&receipt)

	return receipt, nil
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

func (r *ReceiptRepository) Update(receipt *models.Receipt) error {
	return r.db.Omit("created_at").Updates(&receipt).Error
}

func (r *ReceiptRepository) GetByPfr(pfr string) (*models.Receipt, error) {
	var receipt *models.Receipt
	if err := r.db.Where(&models.Receipt{PfrNumber: pfr}).First(&receipt).Error; err != nil {
		return nil, err
	}

	return receipt, nil
}

func (r *ReceiptRepository) GetById(id int) (*models.Receipt, error) {
	var receipt *models.Receipt
	if err := r.db.Preload("ReceiptItems.Category").Preload("Store").First(&receipt, uint(id)).Error; err != nil {
		return nil, err
	}

	return receipt, nil
}
