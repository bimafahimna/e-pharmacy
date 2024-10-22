package usecase

import (
	"context"
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/shopspring/decimal"
)

type PharmacyProductManagementUseCase interface {
	UpdatePharmacy(ctx context.Context) error
	ListMasterProducts(ctx context.Context, pharmacistID int64, params dto.ListProductParams) ([]dto.Product, *dto.Pagination, error)
	ListPharmacyProducts(ctx context.Context, pharmacistID int64, params dto.ListPharmacyProductParams) ([]dto.PharmacyProduct, *dto.Pagination, error)
	GetPharmacyProduct(ctx context.Context, pharmacistID int64, productID int) (*dto.PharmacyProduct, error)
	AddPharmacyProduct(ctx context.Context, pharmacistID int64, req dto.AddPharmacyProductRequest) error
	UpdatePharmacyProduct(ctx context.Context, pharmacistID int64, productID int, req dto.UpdatePharmacyProductRequest) error
	DeletePharmacyProduct(ctx context.Context, pharmacistID int64, productID int) error
}

func (u *pharmacistUseCase) UpdatePharmacy(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *pharmacistUseCase) ListMasterProducts(ctx context.Context, pharmacistID int64, params dto.ListProductParams) ([]dto.Product, *dto.Pagination, error) {
	offset := (params.Page - 1) * params.Limit

	filters := map[string]interface{}{
		"name":                   params.Name,
		"generic_name":           params.GenericName,
		"manufacturer":           params.Manufacturer,
		"product_classification": params.ProductClassification,
		"product_form":           params.ProductForm,
	}

	if params.IsActive != "" {
		isActive, _ := strconv.ParseBool(params.IsActive)
		filters["is_active"] = isActive
	}

	if params.Usage != "" {
		usage, _ := strconv.Atoi(params.Usage)
		filters["usage"] = usage
	}

	pharmacy, err := u.store.Pharmacy().Find(ctx, pharmacistID)
	if err != nil {
		return nil, nil, err
	}

	items, err := u.store.Product().FindAllByPharmacyID(ctx, pharmacy.ID, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.Product{}
	for _, product := range items {
		resp = append(resp, dto.Product{
			ID:                    product.ID,
			Name:                  product.Name,
			GenericName:           product.GenericName,
			Manufacturer:          product.Manufacturer,
			ProductClassification: product.ProductClassification,
			ProductForm:           product.ProductForm,
			IsActive:              product.IsActive,
		})
	}

	totalRecords, err := u.store.Product().CountAllByPharmacyID(ctx, pharmacy.ID, filters)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalRecords + params.Limit - 1) / params.Limit
	currentPage := offset/params.Limit + 1

	pagination := &dto.Pagination{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: currentPage > 1,
		NextPage:     currentPage < totalPages,
	}
	return resp, pagination, nil
}

func (u *pharmacistUseCase) ListPharmacyProducts(ctx context.Context, pharmacistID int64, params dto.ListPharmacyProductParams) ([]dto.PharmacyProduct, *dto.Pagination, error) {
	offset := (params.Page - 1) * params.Limit

	filters := map[string]interface{}{
		"name":                   params.Name,
		"generic_name":           params.GenericName,
		"manufacturer":           params.Manufacturer,
		"product_classification": params.ProductClassification,
		"product_form":           params.ProductForm,
	}

	if params.IsActive != "" {
		isActive, _ := strconv.ParseBool(params.IsActive)
		filters["is_active"] = isActive
	}

	pharmacy, err := u.store.Pharmacy().Find(ctx, pharmacistID)
	if err != nil {
		return nil, nil, err
	}

	items, err := u.store.PharmacyProduct().FindAllByPharmacyID(ctx, pharmacy.ID, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.PharmacyProduct{}
	for _, item := range items {
		resp = append(resp, dto.PharmacyProduct{
			PharmacyID:            item.PharmacyID,
			ProductID:             item.ProductID,
			Stock:                 item.Stock,
			Price:                 item.Price.String(),
			Name:                  item.Name,
			GenericName:           item.GenericName,
			Manufacturer:          item.Manufacturer,
			ProductClassification: item.ProductClassification,
			ProductForm:           item.ProductForm,
			IsActive:              item.IsActive,
		})
	}

	totalRecords, err := u.store.PharmacyProduct().CountAllByPharmacyID(ctx, pharmacy.ID, filters)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalRecords + params.Limit - 1) / params.Limit
	currentPage := offset/params.Limit + 1

	pagination := &dto.Pagination{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: currentPage > 1,
		NextPage:     currentPage < totalPages,
	}
	return resp, pagination, nil
}

func (u *pharmacistUseCase) GetPharmacyProduct(ctx context.Context, pharmacistID int64, productID int) (*dto.PharmacyProduct, error) {
	pharmacy, err := u.store.Pharmacy().Find(ctx, pharmacistID)
	if err != nil {
		return nil, err
	}

	item, err := u.store.PharmacyProduct().Find(ctx, pharmacy.ID, productID)
	if err != nil {
		return nil, err
	}

	resp := &dto.PharmacyProduct{
		PharmacyID:            item.PharmacyID,
		ProductID:             item.ProductID,
		Name:                  item.Name,
		GenericName:           item.GenericName,
		Manufacturer:          item.Manufacturer,
		ProductClassification: item.ProductClassification,
		ProductForm:           item.ProductForm,
		Stock:                 item.Stock,
		IsActive:              item.IsActive,
	}
	return resp, nil
}

func (u *pharmacistUseCase) AddPharmacyProduct(ctx context.Context, pharmacistID int64, req dto.AddPharmacyProductRequest) error {
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return apperror.ErrBadRequest
	}

	return u.store.Atomic(ctx, func(s repository.Store) error {
		pharmacy, err := s.Pharmacy().Find(ctx, pharmacistID)
		if err != nil {
			return err
		}

		exists, err := s.PharmacyProduct().CheckExists(ctx, pharmacy.ID, req.ProductID)
		if err != nil {
			return err
		}
		if exists {
			return apperror.ErrPharmacyProductAlreadyExists
		}

		item := model.PharmacyProduct{
			ProductID:  req.ProductID,
			SoldAmount: 0,
			Price:      price,
			Stock:      req.Stock,
			IsActive:   true,
		}
		if err := s.PharmacyProduct().Insert(ctx, pharmacy.ID, item); err != nil {
			return err
		}
		return s.Product().IncrementUsage(ctx, req.ProductID, 1)
	})
}

func (u *pharmacistUseCase) UpdatePharmacyProduct(ctx context.Context, pharmacistID int64, productID int, req dto.UpdatePharmacyProductRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		pharmacy, err := s.Pharmacy().Find(ctx, pharmacistID)
		if err != nil {
			return err
		}

		updatedToday, err := s.PharmacyProduct().CheckUpdatedToday(ctx, pharmacy.ID, productID)
		if err != nil {
			return err
		}
		if updatedToday {
			return apperror.ErrUpdatedPharmacyProductToday
		}

		isActive, _ := strconv.ParseBool(req.IsActive)

		item := model.PharmacyProduct{
			Stock:    req.Stock,
			IsActive: isActive,
		}
		return s.PharmacyProduct().Update(ctx, pharmacy.ID, productID, item)
	})
}

func (u *pharmacistUseCase) DeletePharmacyProduct(ctx context.Context, pharmacistID int64, productID int) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		// TODO: Check for pharmacy products in cart and remove them first

		pharmacy, err := s.Pharmacy().Find(ctx, pharmacistID)
		if err != nil {
			return err
		}

		sold, err := s.PharmacyProduct().CheckSold(ctx, pharmacy.ID, productID)
		if err != nil {
			return err
		}
		if sold {
			return apperror.ErrPharmacyProductHasBeenBought
		}

		if err := s.PharmacyProduct().Delete(ctx, pharmacy.ID, productID); err != nil {
			return err
		}
		return s.Product().IncrementUsage(ctx, productID, -1)
	})
}
