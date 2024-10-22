package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type orderItemRepository struct {
	db database
}

func NewOrderItemRepository(db database) OrderItemRepository {
	return &orderItemRepository{
		db: db,
	}
}

func (r *orderItemRepository) BulkInsert(ctx context.Context, items []model.OrderItem) error {
	values := []string{}
	args := []interface{}{}

	for i, item := range items {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
		args = append(args, item.OrderID)
		args = append(args, item.PharmacyID)
		args = append(args, item.ProductID)
		args = append(args, item.ImageURL)
		args = append(args, item.Name)
		args = append(args, item.Quantity)
		args = append(args, item.Price)
	}

	query := fmt.Sprintf(`
		INSERT INTO "order_items" (
			"order_id", "pharmacy_id", "product_id", "image_url", "name", "quantity", "price"
		) VALUES %s
	`, strings.Join(values, ","))

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}
