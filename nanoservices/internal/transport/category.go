package transport

import "github.com/steffanturanjanin/receipt-manager/internal/models"

type CategoryResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CategoryTransformer struct{}

func (t CategoryTransformer) TransformSingle(model models.Category) CategoryResponse {
	category := CategoryResponse{}
	category.ID = int(model.ID)
	category.Name = model.Name
	category.Color = model.Color

	return category
}
