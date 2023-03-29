package controllers

import (
	"net/http"

	"github.com/steffanturanjanin/receipt-manager/internal/errors"
	"github.com/steffanturanjanin/receipt-manager/internal/filters"
	"github.com/steffanturanjanin/receipt-manager/internal/services"
)

type StatisticController struct {
	statisticService *services.StatisticService
}

func NewStatisticController(s *services.StatisticService) *StatisticController {
	return &StatisticController{
		statisticService: s,
	}
}

func (c *StatisticController) ListCategoriesStatistic(w http.ResponseWriter, r *http.Request) {
	filters := filters.CategoryStatisticFilters{}
	filters.BuildFromRequest(r)

	categoryStatistics, err := c.statisticService.GetCategoryStatistic(filters)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
	}

	JsonResponse(w, categoryStatistics, http.StatusOK)
}
