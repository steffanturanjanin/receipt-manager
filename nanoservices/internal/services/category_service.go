package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
)

type CategoryService struct {
	categoryRepository repositories.CategoryRepositoryInterface
}

func NewCategoryService(r repositories.CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{
		categoryRepository: r,
	}
}

func (s *CategoryService) GetAll() ([]dto.Category, error) {
	categories := make([]dto.Category, 0)

	categoryModels, err := s.categoryRepository.GetAll()
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

func (s *CategoryService) GetIds() ([]int, error) {
	return s.categoryRepository.GetIds()
}

func (s *CategoryService) GetById(id int) (*dto.Category, error) {
	categoryModel, err := s.categoryRepository.GetById(id)
	if err != nil {
		return nil, err
	}

	return &dto.Category{
		Id:   categoryModel.ID,
		Name: categoryModel.Name,
	}, nil
}
