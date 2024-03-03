package query

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type BaseQueryBuilder struct {
	Query *gorm.DB
}

func NewBaseQueryBuilder(query *gorm.DB) BaseQueryBuilder {
	return BaseQueryBuilder{Query: query}
}

func (qb BaseQueryBuilder) SortableOptions() SortableOptions {
	return []string{}
}

func (qb BaseQueryBuilder) Sort(sortQuery *SortQuery) BaseQueryBuilder {
	for _, option := range qb.SortableOptions() {
		if option == sortQuery.Field {
			qb.Query.Order(fmt.Sprintf("%s %s", sortQuery.Field, sortQuery.Direction))
			return qb
		}
	}

	return qb
}

func (qb BaseQueryBuilder) ExecutePaginatedQuery(destination interface{}, pq PaginationQuery) (*PaginationData, error) {
	// Check if destination is a pointer pointing to a slice
	value := reflect.ValueOf(destination).Elem()
	if value.Kind() != reflect.Slice {
		return nil, errors.New("data must point to a slice")
	}

	// Build pagination meta object
	meta := GetPaginationMeta(qb.Query, pq)

	// Generate offset
	offset := (pq.Page - 1) * pq.Limit

	// Execute query but instead of applying pagination params on Query instance and mutating the Query,
	// perform pagination on cloned query and prevent pollution
	if err := qb.cloneQuery().Limit(pq.Limit).Offset(offset).Find(destination).Error; err != nil {
		return nil, err
	}

	items := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		items[i] = value.Index(i).Interface()
	}

	return &PaginationData{Data: items, Meta: meta}, nil
}

func (qb BaseQueryBuilder) Immutable() BaseQueryBuilder {
	qb.Query.Session(&gorm.Session{})

	return qb
}

func (qb BaseQueryBuilder) GetQuery() *gorm.DB {
	return qb.Query
}

func (qb BaseQueryBuilder) cloneQuery() *gorm.DB {
	return qb.Query.Session(&gorm.Session{})
}
