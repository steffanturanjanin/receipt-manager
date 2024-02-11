package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/models"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
	receipt_fetcher "github.com/steffanturanjanin/receipt-manager/pkg/receipt-fetcher"
)

type StoreService struct {
	storeRepository repositories.StoreRepositoryInterface
}

func NewStoreService(sr repositories.StoreRepositoryInterface) *StoreService {
	return &StoreService{
		storeRepository: sr,
	}
}

func (s *StoreService) GetByTin(tin string) (*dto.Store, error) {
	storeModel, err := s.storeRepository.GetByTin(tin)
	if err != nil {
		return nil, err
	}

	return &dto.Store{
		Tin:          storeModel.Tin,
		Name:         storeModel.Name,
		LocationName: storeModel.LocationName,
		LocationId:   storeModel.LocationId,
		Address:      storeModel.Address,
		City:         storeModel.City,
	}, nil
}

func (s *StoreService) FirstOrCreateFromDto(storeDto receipt_fetcher.Store) (*models.Store, error) {
	return s.storeRepository.FirstOrCreateFromDto(storeDto)
}
