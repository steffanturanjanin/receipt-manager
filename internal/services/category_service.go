package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
)

type CategoryService struct {
	CategoryRepository repositories.CategoryRepositoryInterface
}

func NewCategoryService(r repositories.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{
		CategoryRepository: r,
	}
}

func (s *CategoryService) GetAll() ([]dto.Category, error) {
	categories := make([]dto.Category, 0)

	categoryModels, err := s.CategoryRepository.GetAll()
	if err != nil {
		return nil, err
	}

	for _, categoryModel := range categoryModels {
		categories = append(categories, dto.Category{
			Id:   categoryModel.ID,
			Name: categoryModel.Name,
		})
	}

	return categories, nil
}
