package transport

import "github.com/steffanturanjanin/receipt-manager/internal/models"

type TaxResponse struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Rate       int    `json:"rate"`
}

func (tax TaxResponse) FromModel(model models.Tax) TaxResponse {
	tax.ID = int(model.ID)
	tax.Identifier = model.Identifier
	tax.Name = model.Name
	tax.Rate = model.Rate

	return tax
}
