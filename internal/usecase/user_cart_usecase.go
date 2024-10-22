package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type UserCartUseCase interface {
	GetItems(ctx context.Context, userID int64) ([]dto.CartContent, error)
	AddItem(ctx context.Context, userID int64, req dto.AddCartItemRequest) error
	UpdateItem(ctx context.Context, userId int64, pharmacyID, productID int, req dto.UpdateCartItemRequest) error
	RemoveItem(ctx context.Context, userId int64, pharmacyID, productID int) error
}

func (u *userUseCase) GetItems(ctx context.Context, userID int64) ([]dto.CartContent, error) {
	items, err := u.store.Cart().FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	data := map[string][]dto.CartItem{}
	for _, item := range items {
		cartItem := dto.CartItem{
			PharmacyID:  item.PharmacyID,
			ProductID:   item.ProductID,
			Name:        item.ProductName,
			ImageURL:    item.ImageURL,
			Price:       item.Price,
			Stock:       item.Stock,
			Quantity:    item.Quantity,
			SellingUnit: item.SellingUnit,
			Description: item.Description,
			Weight:      item.Weight.String(),
		}

		data[item.PharmacyName] = append(data[item.PharmacyName], cartItem)
	}

	resp := []dto.CartContent{}
	for name, items := range data {
		content := dto.CartContent{
			PharmacyName: name,
			Items:        items,
		}
		resp = append(resp, content)
	}
	return resp, nil
}

func (u *userUseCase) AddItem(ctx context.Context, userID int64, req dto.AddCartItemRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		item, err := s.PharmacyProduct().Find(ctx, req.PharmacyID, req.ProductID)
		if err != nil {
			return err
		}
		if item.Stock == 0 {
			return apperror.ErrPharmacyProductUnavailable
		}

		exists, err := s.Cart().CheckExists(ctx, userID, req.PharmacyID, req.ProductID)
		if err != nil {
			return err
		}
		if !exists {
			return s.Cart().Insert(ctx, userID, req.PharmacyID, req.ProductID, req.Quantity)
		}
		return s.Cart().Update(ctx, userID, req.PharmacyID, req.ProductID, req.Quantity)
	})
}

func (u *userUseCase) UpdateItem(ctx context.Context, userID int64, pharmacyID, productID int, req dto.UpdateCartItemRequest) error {
	return u.store.Cart().Update(ctx, userID, pharmacyID, productID, req.Quantity)
}

func (u *userUseCase) RemoveItem(ctx context.Context, userID int64, pharmacyID, productID int) error {
	return u.store.Cart().Delete(ctx, userID, pharmacyID, productID)
}
