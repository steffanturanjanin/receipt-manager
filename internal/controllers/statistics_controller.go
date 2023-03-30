package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (c *StatisticController) ListCategoriesStatistics(w http.ResponseWriter, r *http.Request) {
	filters := filters.CategoryStatisticFilters{}
	filters.BuildFromRequest(r)

	categoryStatistics, err := c.statisticService.GetCategoryStatistic(filters)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
	}

	JsonResponse(w, categoryStatistics, http.StatusOK)
}

func (c *StatisticController) ListStoreStatisticsForCategory(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idParam)

	filters := filters.StoreStatisticForCategoryFilters{}
	filters.BuildFromRequest(r)

	result, err := c.statisticService.GetStoreStatisticsForCategory(id, filters)
	if err != nil {
		JsonErrorResponse(w, errors.NewHttpError(err))
		return
	}

	JsonResponse(w, result, http.StatusOK)
}
