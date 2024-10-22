package handler

import "github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"

type PharmacistHandler struct {
	useCase usecase.PharmacistUseCase
}

func NewPharmacistHandler(useCase usecase.PharmacistUseCase) *PharmacistHandler {
	return &PharmacistHandler{
		useCase: useCase,
	}
}
