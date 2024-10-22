package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/hibiken/asynq"
)

type PharmacyOrderManagementUseCase interface {
	ListOrder(ctx context.Context, pharmacistID int64, params dto.ListPharmacyOrderQuery) ([]dto.Order, *dto.Pagination, error)
	SendOrder(ctx context.Context, pharmacistID int64, orderID int, req dto.SendOrderRequest) error
}

func (u *pharmacistUseCase) ListOrder(ctx context.Context, pharmacistID int64, params dto.ListPharmacyOrderQuery) ([]dto.Order, *dto.Pagination, error) {
	offset := (params.Page - 1) * params.Limit
	pharmacy, err := u.store.Pharmacy().Find(ctx, pharmacistID)
	if err != nil {
		return nil, nil, err
	}

	orders, err := u.store.Order().FindAllByPharmacyID(ctx, pharmacy.ID, params.Status, params.SortBy, params.Sort, params.Limit, offset)
	if err != nil {
		return nil, nil, err
	}

	pharmacyOrders := map[int]*dto.Order{}
	for _, order := range orders {
		if _, exists := pharmacyOrders[order.ID]; !exists {
			pharmacyOrders[order.ID] = &dto.Order{
				Address:      order.Address,
				ContactName:  order.ContactName,
				ContactPhone: order.ContactPhone,
				LogisticName: order.LogisticName,
				LogisticCost: order.LogisticCost.String(),
				OrderID:      order.ID,
				OrderItems:   []dto.OrderItem{},
				OrderAmount:  order.Amount.String(),
				CreatedAt:    order.CreatedAt,
				UpdatedAt:    order.UpdatedAt,
			}
		}

		for _, item := range order.Items {
			item := dto.OrderItem{
				PharmacyID: item.PharmacyID,
				ProductID:  item.ProductID,
				ImageURL:   item.ImageURL,
				Name:       item.Name,
				Quantity:   item.Quantity,
				Price:      item.Price.String(),
			}

			pharmacyOrders[order.ID].OrderItems = append(pharmacyOrders[order.ID].OrderItems, item)
		}
	}

	resp := []dto.Order{}
	for _, order := range pharmacyOrders {
		resp = append(resp, *order)
	}

	totalRecords, err := u.store.Order().CountAllByPharmacyID(ctx, pharmacy.ID, params.Status)
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

func (u *pharmacistUseCase) SendOrder(ctx context.Context, pharmacistID int64, orderID int, req dto.SendOrderRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		order, err := s.Order().Find(ctx, orderID)
		if err != nil {
			return err
		}
		if order.Status != "Processed" {
			return apperror.ErrBadRequest
		}

		if err := s.Order().UpdateByOrderID(ctx, orderID, req.Status); err != nil {
			return err
		}

		pharmacy, err := s.Pharmacy().Find(ctx, pharmacistID)
		if err != nil {
			return err
		}

		eda, err := s.PharmacyLogistic().CheckEDA(ctx, order.LogisticID, pharmacy.ID)
		if err != nil {
			if errors.Is(err, apperror.ErrNotFound) {
				return apperror.ErrUnavailableLogistic
			}
			return err
		}

		due := time.Now().Add(time.Duration(eda)).Add(7 * 24 * time.Hour)

		opts := []asynq.Option{
			asynq.Queue(appconst.QueueDefault),
			asynq.MaxRetry(1),
			asynq.ProcessIn(time.Until(due)),
		}
		return u.producer.ProduceConfirmUserOrder(
			ctx,
			&mq.ConfirmUserOrderPayload{
				OrderID: orderID,
				Status:  "Order Confirmed",
			},
			opts...,
		)
	})
}
