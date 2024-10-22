package usecase

import (
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type PharmacistUseCase interface {
	PharmacyProductManagementUseCase
	PharmacyOrderManagementUseCase
}

type pharmacistUseCase struct {
	store    repository.Store
	producer mq.Producer
}

func NewPharmacistUseCase(store repository.Store, producer mq.Producer) PharmacistUseCase {
	return &pharmacistUseCase{
		store:    store,
		producer: producer,
	}
}
