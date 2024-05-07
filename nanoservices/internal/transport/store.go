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

type StoreTransformer struct{}

func (t StoreTransformer) Transform(models []models.Store) []StoreResponse {
	items := make([]StoreResponse, 0)

	for _, model := range models {
		store := t.TransformSingle(model)
		items = append(items, store)
	}

	return items
}

func (t StoreTransformer) TransformSingle(model models.Store) StoreResponse {
	store := StoreResponse{}
	store.ID = int(model.ID)
	store.Tin = model.Tin
	store.Name = model.Name
	store.LocationId = model.LocationId
	store.LocationName = model.LocationName
	store.Address = model.Address
	store.City = model.City

	return store
}
