package transport

type TaxResponse struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Rate       int    `json:"rate"`
}
