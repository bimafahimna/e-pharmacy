package usecase

import (
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/bcrypt"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/cronjob"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type AdminUseCase interface {
	UserManagementUseCase
	PharmacistManagementUseCase
	PartnerManagementUseCase
	PharmacyManagementUseCase
	ProductCategoryManagementUseCase
	ProductManagementUseCase
	AdminOrderUseCase
}

type adminUseCase struct {
	store    repository.Store
	bcrypt   bcrypt.Provider
	producer mq.Producer
	crond    cronjob.Provider
}

func NewAdminUseCase(store repository.Store, bcrypt bcrypt.Provider, producer mq.Producer, crond cronjob.Provider) AdminUseCase {
	return &adminUseCase{
		store:    store,
		bcrypt:   bcrypt,
		producer: producer,
		crond:    crond,
	}
}
