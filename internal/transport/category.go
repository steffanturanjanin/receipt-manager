package transport

import "github.com/steffanturanjanin/receipt-manager/internal/models"

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (category CategoryResponse) FromModel(model models.Category) CategoryResponse {
	category.ID = int(model.ID)
	category.Name = model.Name

	return category
}
