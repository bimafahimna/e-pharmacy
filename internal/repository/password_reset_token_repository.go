package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type passwordResetTokenRepository struct {
	db database
}

func NewPasswordResetTokenRepository(db database) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{
		db: db,
	}
}

func (r *passwordResetTokenRepository) Save(ctx context.Context, userID int64, token string) error {
	query := `INSERT INTO "password_reset_tokens" ("user_id", "token", "expired_at") VALUES ($1, $2, CURRENT_TIMESTAMP + INTERVAL '1 hour')`

	if _, err := r.db.ExecContext(ctx, query, userID, token); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *passwordResetTokenRepository) FindUserID(ctx context.Context, token string) (int64, error) {
	query := `SELECT "user_id" FROM "password_reset_tokens" WHERE "token" = $1 AND "used_at" IS NULL AND "expired_at" > CURRENT_TIMESTAMP`

	var userID int64
	if err := r.db.QueryRowContext(ctx, query, token).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userID, apperror.ErrNotFound
		}
		return userID, apperror.ErrInternalServerError
	}
	return userID, nil
}

func (r *passwordResetTokenRepository) CheckExists(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "password_reset_tokens" where "user_id" = $1 AND "used_at" IS NULL AND "expired_at" > CURRENT_TIMESTAMP)`

	var ok bool
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&ok); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, apperror.ErrInternalServerError
	}
	return ok, nil
}

func (r *passwordResetTokenRepository) Revoke(ctx context.Context, userID int64) error {
	query := `UPDATE "password_reset_tokens" SET "used_at" = CURRENT_TIMESTAMP, "updated_at" = CURRENT_TIMESTAMP WHERE "user_id" = $1`

	if _, err := r.db.ExecContext(ctx, query, userID); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}
