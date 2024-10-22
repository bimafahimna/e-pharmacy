package usecase

import (
	"context"
	"strings"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
	"github.com/shopspring/decimal"
)

type UserHomepageUseCase interface {
	ListBestsellerProduct(ctx context.Context, limit int) ([]dto.Item, error)
	ListRecommendedProduct(ctx context.Context, req dto.ListPopularProductQueries) ([]dto.Item, error)
	ListProduct(ctx context.Context, queries dto.ListProductQueries) ([]dto.Item, *dto.Pagination, error)
}

func (u *userUseCase) ListBestsellerProduct(ctx context.Context, limit int) ([]dto.Item, error) {
	bestSellerID, err := u.store.Product().FindIDBestseller(ctx, limit)
	if err != nil {
		return nil, err
	}

	items, err := u.store.PharmacyProduct().FanOutFanInBestSeller(ctx, bestSellerID, util.MaxWorker(len(bestSellerID)))
	if err != nil {
		return nil, err
	}

	resp := []dto.Item{}
	for _, item := range items {
		resp = append(resp, dto.Item{
			PharmacyID:    item.PharmacyID,
			ProductID:     item.ProductID,
			ProductFormID: item.ProductFormID,
			SellingUnit:   item.SellingUnit,
			Price:         item.Price,
			Stock:         item.Stock,
			Name:          item.Name,
			ImageURL:      item.ImageURL,
		})
	}
	return resp, nil
}

func (u *userUseCase) ListRecommendedProduct(ctx context.Context, queries dto.ListPopularProductQueries) ([]dto.Item, error) {
	latitude, err := decimal.NewFromString(queries.Latitude)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	longitude, err := decimal.NewFromString(queries.Longitude)
	if err != nil {
		return nil, apperror.ErrBadRequest
	}

	minLat, maxLat, minLong, maxLong := util.HaversineBounds(latitude.InexactFloat64(), longitude.InexactFloat64(), 25)

	productIDs, err := u.store.Product().FindIDBestseller(ctx, 26)
	if err != nil {
		return nil, err
	}

	items, err := u.store.PharmacyProduct().FanOutFanInRecommended(ctx, productIDs, minLat, maxLat, minLong, maxLong, util.MaxWorker(len(productIDs)))
	if err != nil {
		return nil, err
	}

	resp := []dto.Item{}
	for _, item := range items {
		resp = append(resp, dto.Item{
			PharmacyID:    item.PharmacyID,
			ProductID:     item.ProductID,
			ProductFormID: item.ProductFormID,
			SellingUnit:   item.SellingUnit,
			Price:         item.Price,
			Stock:         item.Stock,
			Name:          item.Name,
			ImageURL:      item.ImageURL,
		})
	}
	return resp, nil
}

func (u *userUseCase) ListProduct(ctx context.Context, queries dto.ListProductQueries) ([]dto.Item, *dto.Pagination, error) {
	offset := util.ToOffset(queries.Page, queries.Limit)
	ensureSearch := util.RemoveSymbols(queries.Search)
	searchTrim := strings.TrimSpace(ensureSearch)
	search := strings.ReplaceAll(searchTrim, " ", " & of & ")
	productId, err := u.store.Product().SearchID(ctx, search, offset, queries.Limit+1)
	if err != nil {
		return nil, nil, err
	}
	nextPage := len(productId) > queries.Limit
	if nextPage {
		productId = productId[:queries.Limit]
	}

	var (
		items []model.Item
	)
	if queries.Latitude == "0" || queries.Longitude == "0" {
		items, err = u.store.PharmacyProduct().FanOutFanInBestSeller(ctx, productId, util.MaxWorker(len(productId)))
		if err != nil {
			return nil, nil, err
		}
	} else {
		latitude, err := decimal.NewFromString(queries.Latitude)
		if err != nil {
			return nil, nil, apperror.ErrBadRequest
		}
		longitude, err := decimal.NewFromString(queries.Longitude)
		if err != nil {
			return nil, nil, apperror.ErrBadRequest
		}
		minLat, maxLat, minLong, maxLong := util.HaversineBounds(latitude.InexactFloat64(), longitude.InexactFloat64(), 25)
		items, err = u.store.PharmacyProduct().FanOutFanInRecommended(ctx, productId, minLat, maxLat, minLong, maxLong, util.MaxWorker(len(productId)))
		if err != nil {
			return nil, nil, err
		}
	}

	resp := []dto.Item{}
	for _, item := range items {
		resp = append(resp, dto.Item{
			PharmacyID:    item.PharmacyID,
			ProductID:     item.ProductID,
			ProductFormID: item.ProductFormID,
			SellingUnit:   item.SellingUnit,
			Price:         item.Price,
			Stock:         item.Stock,
			Name:          item.Name,
			ImageURL:      item.ImageURL,
		})
	}
	pagination := &dto.Pagination{
		CurrentPage: queries.Page,
		NextPage:    nextPage,
	}
	return resp, pagination, nil
}
