package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type verifyTokenRepository struct {
	db database
}

func NewVerifyTokenRepository(db database) VerificationTokenRepository {
	return &verifyTokenRepository{
		db: db,
	}
}

func (r *verifyTokenRepository) Save(ctx context.Context, userID int64, token string) error {
	query := `
		INSERT INTO "verification_tokens" ("user_id", "token", "expired_at")
		VALUES ($1, $2, CURRENT_TIMESTAMP + INTERVAL '1 hour')
	`
	args := []interface{}{userID, token}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *verifyTokenRepository) Find(ctx context.Context, token string) (int64, error) {
	query := `
		SELECT "user_id" FROM "verification_tokens"
		WHERE "token" = $1 AND "used_at" IS NULL AND "expired_at" > CURRENT_TIMESTAMP
	`

	var userID int64
	if err := r.db.QueryRowContext(ctx, query, token).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperror.ErrNotFound
		}
		return 0, apperror.ErrInternalServerError
	}
	return userID, nil
}

func (r *verifyTokenRepository) Revoke(ctx context.Context, userID int64) error {
	query := `
		UPDATE "verification_tokens"
		SET "used_at" = CURRENT_TIMESTAMP, "updated_at" = CURRENT_TIMESTAMP
		WHERE "user_id" = $1
	`

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *verifyTokenRepository) CheckExists(ctx context.Context, userID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "verification_tokens" 
			WHERE "user_id" = $1 AND "used_at" IS NULL AND "expired_at" > CURRENT_TIMESTAMP
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
