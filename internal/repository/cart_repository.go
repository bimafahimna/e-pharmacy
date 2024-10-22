package repository

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository/postgres"
)

type cartRepository struct {
	db database
}

func NewCartRepository(db database) CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) FindAll(ctx context.Context, userID int64) ([]model.CartItem, error) {
	query := `
		SELECT 
			pp."pharmacy_id", pp."product_id", ph."name", 
			p."name", p."image_url", p."description", 
			p."selling_unit", pp."price", pp."stock", c."quantity", p."weight"
		FROM "carts" c
		JOIN "pharmacy_products" pp ON
			pp."pharmacy_id" = c."pharmacy_id" AND pp."product_id" = c."product_id"
		JOIN "pharmacies" ph ON ph."id" = pp."pharmacy_id"
		JOIN "products" p ON p."id" = pp."product_id"
		WHERE c."user_id" = $1
	`
	args := []interface{}{userID}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	items := []model.CartItem{}
	for rows.Next() {
		var item model.CartItem
		dest := []interface{}{
			&item.PharmacyID, &item.ProductID, &item.PharmacyName,
			&item.ProductName, &item.ImageURL, &item.Description,
			&item.SellingUnit, &item.Price, &item.Stock, &item.Quantity, &item.Weight,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return items, nil
}

func (r *cartRepository) Insert(ctx context.Context, userID int64, pharmacyID, productID, quantity int) error {
	query := `
		INSERT INTO "carts" ("user_id", "pharmacy_id", "product_id", "quantity") 
		VALUES ($1, $2, $3, $4)
	`
	args := []interface{}{userID, pharmacyID, productID, quantity}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		if postgres.IsForeignKeyViolation(err) {
			return apperror.ErrBadRequest
		}
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *cartRepository) Update(ctx context.Context, userID int64, pharmacyID, productID, quantity int) error {
	query := `
		UPDATE "carts" SET "quantity" = $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "user_id" = $2 AND "pharmacy_id" = $3 AND "product_id" = $4
	`
	args := []interface{}{quantity, userID, pharmacyID, productID}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNoChangesMade
	}
	return nil
}

func (r *cartRepository) Delete(ctx context.Context, userID int64, pharmacyID, productID int) error {
	query := `
		DELETE FROM "carts" WHERE "user_id" = $1 AND "pharmacy_id" = $2 AND "product_id" = $3
	`
	args := []interface{}{userID, pharmacyID, productID}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNoChangesMade
	}
	return nil
}

func (r *cartRepository) DeleteOrdered(ctx context.Context, userID int64, pharmacies []int, products []int) error {
	query := `
		DELETE FROM "carts"
		WHERE "user_id" = $1
			AND "pharmacy_id" = ANY($2::int[])
			AND "product_id" = ANY($3::int[])
	`
	args := []interface{}{userID, pharmacies, products}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *cartRepository) CheckExists(ctx context.Context, userID int64, pharmacyID, productID int) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM "carts"
		WHERE "user_id" = $1 AND "pharmacy_id" = $2 AND "product_id" = $3
	)`
	args := []interface{}{userID, pharmacyID, productID}

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
