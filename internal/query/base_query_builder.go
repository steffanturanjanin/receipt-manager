package query

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BaseQueryBuilder struct {
	Query *gorm.DB
}

func NewBaseQueryBuilder(query *gorm.DB) BaseQueryBuilder {
	return BaseQueryBuilder{Query: query}
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

func (qb BaseQueryBuilder) SortableOptions() SortableOptions {
	return []string{}
}

func (qb BaseQueryBuilder) GetQuery() *gorm.DB {
	return qb.Query
}

func (qb BaseQueryBuilder) ExecutePaginatedQuery(destination interface{}, pq *PaginationQuery) (*PaginationData, error) {
	items, ok := destination.([]interface{})
	if !ok {
		return nil, errors.New("parameter destination must be slice")
	}

	if err := qb.Query.Limit(pq.Limit).Offset(pq.Page).Find(&destination).Error; err != nil {
		return nil, err
	}

	meta := GetPaginationMeta(qb.Query, *pq)

	return &PaginationData{Data: items, Meta: meta}, nil
}
