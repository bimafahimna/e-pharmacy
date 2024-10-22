package usecase

import (
	"context"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type ProductCategoryManagementUseCase interface {
	ListProductCategories(ctx context.Context) ([]dto.ProductCategory, error)
	AddProductCategory(ctx context.Context, req dto.AddProductCategoryRequest) error
	UpdateProductCategory(ctx context.Context, id int, req dto.UpdateProductCategoryRequest) error
	RemoveProductCategory(ctx context.Context, id int) error
}

func (u *adminUseCase) ListProductCategories(ctx context.Context) ([]dto.ProductCategory, error) {
	categories, err := u.store.ProductCategory().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	resp := []dto.ProductCategory{}
	for _, category := range categories {
		resp = append(resp, dto.ProductCategory{
			ID:   category.ID,
			Name: category.Name,
		})
	}
	return resp, nil
}

func (u *adminUseCase) AddProductCategory(ctx context.Context, req dto.AddProductCategoryRequest) error {
	category := model.ProductCategory{
		Name: req.Name,
	}
	if err := u.store.ProductCategory().Create(ctx, category); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return apperror.ErrProductCategoryMustBeUnique
		}
		return err
	}
	return nil
}

func (u *adminUseCase) UpdateProductCategory(ctx context.Context, id int, req dto.UpdateProductCategoryRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		exists, err := s.ProductCategory().CheckExists(ctx, req.Name)
		if err != nil {
			return err
		}
		if exists {
			return apperror.ErrProductCategoryMustBeUnique
		}

		return s.ProductCategory().Update(ctx, id, model.ProductCategory{
			Name: req.Name,
		})
	})
}

func (u *adminUseCase) RemoveProductCategory(ctx context.Context, id int) error {
	return u.store.ProductCategory().Delete(ctx, id)
}
