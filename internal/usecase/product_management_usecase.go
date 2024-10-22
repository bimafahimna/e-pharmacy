package usecase

import (
	"context"
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type ProductManagementUseCase interface {
	AddProduct(ctx context.Context, req dto.AddProductRequest) error
	ListProduct(ctx context.Context, params dto.ListProductParams) ([]dto.Product, *dto.Pagination, error)
	UpdateProduct(ctx context.Context) error
	RemoveProduct(ctx context.Context) error
	SearchProducts(ctx context.Context) error
	ListManufacturer(ctx context.Context, params dto.ListProductDetailParams) ([]dto.ProductDetail, error)
	ListProductClassification(ctx context.Context, params dto.ListProductDetailParams) ([]dto.ProductDetail, error)
	ListProductForm(ctx context.Context, params dto.ListProductDetailParams) ([]dto.ProductDetail, error)
}

func (u *adminUseCase) AddProduct(ctx context.Context, req dto.AddProductRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		productRepository := s.Product()
		productClassificationRepository := s.ProductClassification()
		productFormRepository := s.ProductForm()
		manufacturerRepository := s.Manufacturer()

		ok, err := productClassificationRepository.CheckExists(ctx, req.ProductClassificationID)
		if err != nil {
			return err
		}
		if !ok {
			return apperror.ErrInvalidProductClassification
		}

		ok, err = productFormRepository.CheckExists(ctx, req.ProductFormID)
		if err != nil {
			return err
		}
		if !ok {
			return apperror.ErrInvalidProductForm
		}

		ok, err = manufacturerRepository.CheckExists(ctx, req.ManufacturerID)
		if err != nil {
			return err
		}
		if !ok {
			return apperror.ErrInvalidManufacturer
		}

		ok, err = productRepository.CheckExists(ctx, req.Name, req.GenericName, req.ManufacturerID)
		if err != nil {
			return err
		}
		if ok {
			return apperror.ErrProductNameHasBeenRegistered
		}

		product := dto.ConvertToProductModel(req)

		err = productRepository.InsertItem(ctx, *product)

		return err
	})
}

func (u *adminUseCase) ListProduct(ctx context.Context, params dto.ListProductParams) ([]dto.Product, *dto.Pagination, error) {
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

	products, err := u.store.Product().FindAll(ctx, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.Product{}
	for _, product := range products {
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

	totalRecords, err := u.store.Product().CountAll(ctx, filters)
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

func (u *adminUseCase) UpdateProduct(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *adminUseCase) RemoveProduct(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *adminUseCase) SearchProducts(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (u *adminUseCase) ListManufacturer(ctx context.Context, params dto.ListProductDetailParams) ([]dto.ProductDetail, error) {
	manufacturers, err := u.store.Manufacturer().FindAll(ctx, params.Name)
	if err != nil {
		return nil, err
	}

	resp := []dto.ProductDetail{}
	for _, manufacturer := range manufacturers {
		resp = append(resp, dto.ProductDetail{
			ID:   manufacturer.ID,
			Name: manufacturer.Name,
		})
	}
	return resp, nil
}

func (u *adminUseCase) ListProductClassification(ctx context.Context, params dto.ListProductDetailParams) ([]dto.ProductDetail, error) {
	classifications, err := u.store.ProductClassification().FindAll(ctx, params.Name)
	if err != nil {
		return nil, err
	}

	resp := []dto.ProductDetail{}
	for _, classification := range classifications {
		resp = append(resp, dto.ProductDetail{
			ID:   classification.ID,
			Name: classification.Name,
		})
	}
	return resp, nil
}

func (u *adminUseCase) ListProductForm(ctx context.Context, params dto.ListProductDetailParams) ([]dto.ProductDetail, error) {
	forms, err := u.store.ProductForm().FindAll(ctx, params.Name)
	if err != nil {
		return nil, err
	}

	resp := []dto.ProductDetail{}
	for _, form := range forms {
		resp = append(resp, dto.ProductDetail{
			ID:   form.ID,
			Name: form.Name,
		})
	}
	return resp, nil
}
