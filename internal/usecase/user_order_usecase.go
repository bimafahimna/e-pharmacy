package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/hibiken/asynq"
	"github.com/shopspring/decimal"
)

type UserOrderUseCase interface {
	CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) error
	ListUnpaidOrder(ctx context.Context, userID int64) ([]dto.UnpaidOrder, error)
	ListOrder(ctx context.Context, userID int64, status string) ([]dto.Order, error)
	UploadPaymentProof(ctx context.Context, paymentID int, imageURL string) error
	UpdateOrderStatus(ctx context.Context, orderID int, userID int64, status string) error
}

func (u *userUseCase) CreateOrder(ctx context.Context, userID int64, req dto.CreateOrderRequest) error {
	amount, err := decimal.NewFromString(req.PaymentAmount)
	if err != nil {
		return apperror.ErrBadRequest
	}
	productOccurences := map[int]int{}
	for i := 0; i < len(req.Orders); i++ {
		for j := 0; j < len(req.Orders[i].OrderItems); j++ {
			_, ok := productOccurences[req.Orders[i].OrderItems[j].ProductID]
			if ok {
				return apperror.ErrBadRequest
			} else {
				productOccurences[req.Orders[i].OrderItems[j].ProductID] += 1
			}
		}
	}

	return u.store.Atomic(ctx, func(s repository.Store) error {
		payment := model.Payment{
			PaymentMethod: req.PaymentMethod,
			Amount:        amount,
		}
		paymentID, err := s.Payment().Create(ctx, payment)
		if err != nil {
			return err
		}

		pharmacies := []int{}
		products := []int{}

		for _, order := range req.Orders {
			logisticCost, err := decimal.NewFromString(order.LogisticCost)
			if err != nil {
				return apperror.ErrBadRequest
			}

			orderAmount, err := decimal.NewFromString(order.OrderAmount)
			if err != nil {
				return apperror.ErrBadRequest
			}

			newOrder := model.Order{
				UserID:       userID,
				PharmacyID:   order.PharmacyID,
				PaymentID:    paymentID,
				Status:       "Waiting for payment",
				Address:      order.Address,
				PharmacyName: order.PharmacyName,
				ContactName:  order.ContactName,
				ContactPhone: order.ContactPhone,
				LogisticID:   order.LogisticID,
				LogisticCost: logisticCost,
				Amount:       orderAmount,
			}
			orderID, err := s.Order().Create(ctx, paymentID, newOrder)
			if err != nil {
				return err
			}

			items := []model.OrderItem{}
			for _, item := range order.OrderItems {
				price, err := decimal.NewFromString(item.Price)
				if err != nil {
					return apperror.ErrBadRequest
				}

				items = append(items, model.OrderItem{
					OrderID:    orderID,
					PharmacyID: item.PharmacyID,
					ProductID:  item.ProductID,
					ImageURL:   item.ImageURL,
					Name:       item.Name,
					Quantity:   item.Quantity,
					Price:      price,
				})
				if err := s.PharmacyProduct().UpdateStock(ctx, item.PharmacyID, item.ProductID, -item.Quantity); err != nil {
					return err
				}

				pharmacies = append(pharmacies, item.PharmacyID)
				products = append(products, item.ProductID)

			}

			if err := s.OrderItem().BulkInsert(ctx, items); err != nil {
				return err
			}
		}

		return s.Cart().DeleteOrdered(ctx, userID, pharmacies, products)
	})
}

func (u *userUseCase) ListUnpaidOrder(ctx context.Context, userID int64) ([]dto.UnpaidOrder, error) {
	orders, err := u.store.Order().FindAllUnpaid(ctx, userID)
	if err != nil {
		return nil, err
	}

	payments := map[int]*dto.UnpaidOrder{}
	for _, order := range orders {
		orderItem := dto.OrderItem{
			PharmacyID: order.PharmacyID,
			ProductID:  order.ProductID,
			ImageURL:   order.ImageURL,
			Name:       order.Name,
			Quantity:   order.Quantity,
			Price:      order.Price.String(),
		}

		if _, ok := payments[order.PaymentID]; !ok {
			payments[order.PaymentID] = &dto.UnpaidOrder{
				PaymentID:     order.PaymentID,
				PaymentAmount: order.PaymentAmount,
				Orders:        []dto.Order{},
			}
		}

		var existingOrder *dto.Order
		for i, o := range payments[order.PaymentID].Orders {
			if o.OrderID == order.ID {
				existingOrder = &payments[order.PaymentID].Orders[i]
				break
			}
		}
		if existingOrder == nil {
			existingOrder = &dto.Order{
				Address:         order.Address,
				OrderID:         order.ID,
				PharmacyName:    order.PharmacyName,
				OrderItems:      []dto.OrderItem{orderItem},
				ContactName:     order.ContactName,
				ContactPhone:    order.ContactPhone,
				LogisticName:    order.LogisticName,
				LogisticService: order.LogisticService,
				LogisticCost:    order.LogisticCost.String(),
				OrderAmount:     order.OrderAmount.String(),
			}
			payments[order.PaymentID].Orders = append(payments[order.PaymentID].Orders, *existingOrder)
			continue
		}

		existingOrder.OrderItems = append(existingOrder.OrderItems, orderItem)
	}

	resp := []dto.UnpaidOrder{}
	for _, unpaidOrder := range payments {
		resp = append(resp, *unpaidOrder)
	}
	return resp, nil
}

func (u *userUseCase) ListOrder(ctx context.Context, userID int64, status string) ([]dto.Order, error) {
	orders, err := u.store.Order().FindAll(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	pharmacyOrders := map[int]*dto.Order{}
	for _, order := range orders {
		if _, exists := pharmacyOrders[order.ID]; !exists {
			pharmacyOrders[order.ID] = &dto.Order{
				Status:          order.Status,
				Address:         order.Address,
				ContactName:     order.ContactName,
				ContactPhone:    order.ContactPhone,
				PharmacyID:      order.PharmacyID,
				PharmacyName:    order.PharmacyName,
				LogisticName:    order.LogisticName,
				LogisticService: order.LogisticService,
				LogisticCost:    order.LogisticCost.String(),
				OrderID:         order.ID,
				OrderItems:      []dto.OrderItem{},
				OrderAmount:     order.Amount.String(),
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
	return resp, nil
}

func (u *userUseCase) UploadPaymentProof(ctx context.Context, paymentID int, imageURL string) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		if err := s.Payment().Save(ctx, paymentID, imageURL); err != nil {
			if errors.Is(err, apperror.ErrNotFound) {
				return apperror.ErrPaymentNotFound
			}
			return err
		}

		opts := []asynq.Option{
			asynq.Queue(appconst.QueueCritical),
			asynq.MaxRetry(1),
			asynq.ProcessIn(1 * time.Minute),
		}
		return u.producer.ProduceApprovePaymentProof(
			ctx,
			&mq.ApprovePaymentProofPayload{
				PaymentID: paymentID,
				Status:    "Processed",
			},
			opts...,
		)
	})
}

func (u *userUseCase) UpdateOrderStatus(ctx context.Context, orderID int, userID int64, status string) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		exists, err := s.Order().CheckExists(ctx, userID, int64(orderID))
		if err != nil {
			return err
		}

		if !exists {
			return apperror.ErrBadRequest
		}

		return s.Order().UpdateByOrderID(ctx, orderID, status)
	})
}
