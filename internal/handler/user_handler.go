package handler

import (
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/cache"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"
)

type UserHandler struct {
	useCase usecase.UserUseCase
	cache   cache.Provider
	config  *config.Config
}

func NewUserHandler(useCase usecase.UserUseCase, cache cache.Provider, config *config.Config) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		cache:   cache,
		config:  config,
	}
}
