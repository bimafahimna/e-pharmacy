package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type AdminOrderUseCase interface {
	UpdateOrderStatus(ctx context.Context, paymentID int, status string) error
}

func (u *adminUseCase) UpdateOrderStatus(ctx context.Context, paymentID int, status string) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		return s.Order().UpdateByPaymentID(ctx, paymentID, status)
	})
}
