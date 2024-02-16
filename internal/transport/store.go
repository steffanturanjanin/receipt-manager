package transport

type StoreResponse struct {
	ID           int    `json:"id"`
	Tin          string `json:"tin"`
	Name         string `json:"name"`
	LocationId   string `json:"locationId"`
	LocationName string `json:"locationName"`
	Address      string `json:"address"`
	City         string `json:"city"`
}
