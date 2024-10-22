package usecase

import (
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/bcrypt"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/jwt"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"golang.org/x/oauth2"
)

type UserUseCase interface {
	UserRegistrationUseCase
	UserLoginUseCase
	UserPasswordResetUseCase
	UserHomepageUseCase
	UserProfileUseCase
	UserCartUseCase
	UserOrderUseCase
	UserProductDetailUseCase
}

type userUseCase struct {
	store    repository.Store
	bcrypt   bcrypt.Provider
	jwt      jwt.Provider
	producer mq.Producer
	config   oauth2.Config
}

func NewUserUseCase(store repository.Store, bcrypt bcrypt.Provider, jwt jwt.Provider, producer mq.Producer, config oauth2.Config) UserUseCase {
	return &userUseCase{
		store:    store,
		bcrypt:   bcrypt,
		jwt:      jwt,
		producer: producer,
		config:   config,
	}
}
