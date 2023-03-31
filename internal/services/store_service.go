package services

import (
	"github.com/steffanturanjanin/receipt-manager/internal/dto"
	"github.com/steffanturanjanin/receipt-manager/internal/repositories"
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
