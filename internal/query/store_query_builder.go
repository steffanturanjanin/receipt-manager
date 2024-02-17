package query

import (
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"gorm.io/gorm"
)

type StoreQueryBuilder struct {
	BaseQueryBuilder
}

func NewStoreQueryBuilder(query *gorm.DB) StoreQueryBuilder {
	query = query.Model(&models.Store{})
	baseQueryBuilder := BaseQueryBuilder{Query: query}

	return StoreQueryBuilder{
		BaseQueryBuilder: baseQueryBuilder,
	}
}
