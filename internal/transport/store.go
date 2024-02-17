package transport

import "github.com/steffanturanjanin/receipt-manager/internal/models"

type StoreResponse struct {
	ID           int    `json:"id"`
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	LocationId   string `json:"locationId"`
	LocationName string `json:"locationName"`
	Address      string `json:"address"`
	City         string `json:"city"`
}

type StoresResponse struct {
	Items []StoreResponse
}

func (store StoreResponse) FromModel(model models.Store) StoreResponse {
	store.ID = int(model.ID)
	store.Tin = model.Tin
	store.LocationId = model.LocationId
	store.LocationName = model.LocationName
	store.Address = model.Address
	store.City = model.City

	return store
}

func (stores StoresResponse) FromModels(models []models.Store) StoresResponse {
	for _, model := range models {
		store := StoreResponse{}
		store = store.FromModel(model)
		stores.Items = append(stores.Items, store)
	}

	return stores
}
