package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type LogisticUseCase interface {
	ListLogistics(ctx context.Context) ([]dto.Logistic, error)
}

type logisticUseCase struct {
	store repository.Store
}

func NewLogisticUseCase(store repository.Store) LogisticUseCase {
	return &logisticUseCase{
		store: store,
	}
}

func (u *logisticUseCase) ListLogistics(ctx context.Context) ([]dto.Logistic, error) {
	logistics, err := u.store.LogisticRepository().FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.Logistic, len(logistics))
	for i, logistic := range logistics {
		res[i] = dto.Logistic{
			ID:      logistic.ID,
			Name:    logistic.Name,
			LogoUrl: logistic.LogoUrl,
			Service: logistic.Service,
		}
	}
	return res, nil
}
