package handler

import "github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"

type AdminHandler struct {
	useCase usecase.AdminUseCase
}

func NewAdminHandler(adminUseCase usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		useCase: adminUseCase,
	}
}
