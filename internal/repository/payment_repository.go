package repository

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type paymentRepository struct {
	db database
}

func NewPaymentRepository(db database) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment model.Payment) (int, error) {
	query := `
        INSERT INTO "payments" ("payment_method", "image_url", "amount")
        VALUES ($1, $2, $3) RETURNING "id"
    `
	args := []interface{}{payment.PaymentMethod, payment.ImageURL, payment.Amount}

	var id int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return id, nil
}

func (r *paymentRepository) Save(ctx context.Context, id int, imageURL string) error {
	query := `
		UPDATE "payments"
		SET "image_url" = $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $2
	`
	args := []interface{}{imageURL, id}

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
