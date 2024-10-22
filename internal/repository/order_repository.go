package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type orderRepository struct {
	db database
}

func NewOrderRepository(db database) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) Create(ctx context.Context, paymentID int, order model.Order) (int, error) {
	query := `
		INSERT INTO "orders" (
			"user_id", "pharmacy_id", "payment_id", "status", "address", "pharmacy_name",
			"contact_name", "contact_phone", "logistic_id", "logistic_cost", "amount"
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING "id"
	`
	args := []interface{}{
		order.UserID, order.PharmacyID, paymentID, order.Status, order.Address, order.PharmacyName,
		order.ContactName, order.ContactPhone, order.LogisticID, order.LogisticCost, order.Amount,
	}

	var id int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return id, nil
}

func (r *orderRepository) FindAllUnpaid(ctx context.Context, userID int64) ([]model.UnpaidOrder, error) {
	query := `
		SELECT 
			o."id", o."pharmacy_id", oi."product_id",  o."payment_id", o."pharmacy_name", o."address", oi."image_url", oi."name", oi."quantity",
			oi."price", o."contact_name", o."contact_phone", l."name", l."service", o."logistic_cost", o."amount", p."amount"
		FROM "orders" o
		JOIN "order_items" oi ON oi."order_id" = o."id"
		JOIN "payments" p ON p."id" = o."payment_id"
		JOIN "logistics" l ON l."id" = o."logistic_id"
		WHERE o."status" ILIKE 'Waiting for payment' AND o."user_id" = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	orders := []model.UnpaidOrder{}
	for rows.Next() {
		var order model.UnpaidOrder
		dest := []interface{}{
			&order.ID, &order.PharmacyID, &order.ProductID, &order.PaymentID, &order.PharmacyName, &order.Address, &order.ImageURL, &order.Name, &order.Quantity,
			&order.Price, &order.ContactName, &order.ContactPhone, &order.LogisticName, &order.LogisticService, &order.LogisticCost, &order.OrderAmount, &order.PaymentAmount,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return orders, nil
}

func (r *orderRepository) FindAll(ctx context.Context, userID int64, status string) ([]model.Order, error) {
	query := `
		SELECT 
			o."id", o."pharmacy_id", o."status", o."pharmacy_name", o."address", 
			o."contact_name", o."contact_phone", l."name", l."service",
			o."logistic_cost", o."amount", oi."pharmacy_id", oi."product_id", 
			oi."image_url", oi."name", oi."quantity", oi."price" 
		FROM "orders" o
		JOIN "order_items" oi ON oi."order_id" = o."id"
		JOIN "logistics" l ON l."id" = o."logistic_id"
		WHERE o."user_id" = $1 AND o."status" ILIKE $2
	`
	args := []interface{}{userID, status}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	table := map[int]*model.Order{}
	for rows.Next() {
		var order model.Order
		var item model.OrderItem
		dest := []interface{}{
			&order.ID, &order.PharmacyID, &order.Status, &order.PharmacyName, &order.Address,
			&order.ContactName, &order.ContactPhone, &order.LogisticName, &order.LogisticService,
			&order.LogisticCost, &order.Amount, &item.PharmacyID, &item.ProductID,
			&item.ImageURL, &item.Name, &item.Quantity, &item.Price,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}

		if _, exists := table[order.ID]; !exists {
			order.Items = []model.OrderItem{}
			table[order.ID] = &order
		}

		table[order.ID].Items = append(table[order.ID].Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	orders := []model.Order{}
	for _, order := range table {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (r *orderRepository) FindAllByPharmacyID(ctx context.Context, pharmacyID int, status, sortBy, sort string, limit, offset int) ([]model.Order, error) {
	query := `
		select 
			o."id",
			o."user_id",
			o."address", 
			o."contact_name",
			o."contact_phone",
			l."name",
			l."service",
			o."logistic_cost",
			o."amount",
			o."created_at",
			o."updated_at",
			oi."pharmacy_id",
			oi."product_id", 
			oi."image_url",
			oi."name",
			oi."quantity",
			oi."price"
		from
			(
			select
				*
			from
				orders o1
			where 
				o1."pharmacy_id" = $1
				and o1."status" ilike $2
			limit $3 offset $4) o	 
		join "order_items" oi on
			oi."order_id" = o."id"
		join "logistics" l on
			l."id" = o."logistic_id"
	`

	args := []interface{}{pharmacyID, status, limit, offset}
	query += fmt.Sprintf(`
		order by o."%s" %s  	
	`, sortBy, sort)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	table := map[int]*model.Order{}
	for rows.Next() {
		var order model.Order
		var item model.OrderItem
		dest := []interface{}{
			&order.ID,
			&order.UserID,
			&order.Address,
			&order.ContactName,
			&order.ContactPhone,
			&order.LogisticName,
			&order.LogisticService,
			&order.LogisticCost,
			&order.Amount,
			&order.CreatedAt,
			&order.UpdatedAt,
			&item.PharmacyID,
			&item.ProductID,
			&item.ImageURL,
			&item.Name,
			&item.Quantity,
			&item.Price,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}

		if _, exists := table[order.ID]; !exists {
			order.Items = []model.OrderItem{}
			table[order.ID] = &order
		}

		table[order.ID].Items = append(table[order.ID].Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	orders := []model.Order{}
	for _, order := range table {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (r *orderRepository) CountAllByPharmacyID(ctx context.Context, pharmacyID int, status string) (int, error) {
	query := `
		select
			count(o.*)
		from
			orders o
		where
			o.pharmacy_id = $1
			and o.status ilike $2;
	`
	args := []interface{}{pharmacyID, status}

	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return total, nil
}

func (r *orderRepository) Find(ctx context.Context, id int) (*model.Order, error) {
	query := `
		SELECT "id", "user_id", "pharmacy_id", "payment_id", "logistic_id", "status"
		FROM "orders" WHERE "id" = $1
	`
	var order model.Order
	dest := []interface{}{&order.ID, &order.UserID, &order.PharmacyID, &order.PaymentID, &order.LogisticID, &order.Status}

	if err := r.db.QueryRowContext(ctx, query, id).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &order, nil
}

func (r *orderRepository) UpdateByOrderID(ctx context.Context, orderID int, status string) error {
	query := `
		UPDATE "orders"
		SET "status" = $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $2
	`
	args := []interface{}{status, orderID}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *orderRepository) UpdateByPaymentID(ctx context.Context, paymentID int, status string) error {
	query := `
		UPDATE "orders"
		SET "status" = $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "payment_id" = $2 
	`
	args := []interface{}{status, paymentID}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *orderRepository) CheckExists(ctx context.Context, id int64, orderID int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "users" u JOIN orders o ON u.id = o.user_id WHERE u."id" = $1 AND o.id = $2 AND o.status = 'Sent');`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id, orderID).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
