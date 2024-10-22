package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type userRepository struct {
	db database
}

func NewUserRepository(db database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user model.User) (int64, error) {
	query := `
		INSERT INTO "users" ("role", "email", "password_hash", "is_verified") 
		VALUES ($1, $2, $3, $4) RETURNING "id"
	`
	args := []interface{}{user.Role, user.Email, user.PasswordHash, user.IsVerified}

	var id int64
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return id, nil
}

func (r *userRepository) FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.User, error) {
	query := `SELECT "id", "role", "email", "is_verified", "created_at", "updated_at" FROM "users" WHERE 1=1`

	args := []interface{}{}

	for filter, value := range filters {
		if value == nil {
			continue
		}
		switch filter {
		case "role", "email":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_verified":
			query += fmt.Sprintf(` AND "%s" = $%d`, filter, len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	query += fmt.Sprintf(`
		ORDER BY "%s" %s
		OFFSET $%d
		LIMIT $%d
	`, sortBy, sort, len(args)+1, len(args)+2)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Role, &u.Email, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, apperror.ErrInternalServerError
		}

		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return users, nil
}

func (r *userRepository) CountAll(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `SELECT COUNT(*) FROM "users" WHERE 1=1`

	args := []interface{}{}
	for filter, value := range filters {
		if value == nil {
			continue
		}
		switch filter {
		case "role", "email":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_verified":
			query += fmt.Sprintf(` AND "%s" = $%d`, filter, len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT "role", "email", "password_hash", "is_verified"
		FROM "users" WHERE "id" = $1
	`

	var user model.User
	dest := []interface{}{&user.Role, &user.Email, &user.PasswordHash, &user.IsVerified}

	if err := r.db.QueryRowContext(ctx, query, id).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrEmailNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT "id", "role", "password_hash", "is_verified", "deleted_at"
		FROM "users" WHERE "email" ILIKE $1
	`

	var user model.User
	dest := []interface{}{&user.ID, &user.Role, &user.PasswordHash, &user.IsVerified, &user.DeletedAt}

	if err := r.db.QueryRowContext(ctx, query, email).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrEmailNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &user, nil
}

func (r *userRepository) Verify(ctx context.Context, id int64) error {
	query := `
		UPDATE "users"
		SET "is_verified" = true, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $1
	`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `
		UPDATE "users"
		SET "password_hash" = $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $2
	`

	if _, err := r.db.ExecContext(ctx, query, passwordHash, id); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE "users" 
		SET "deleted_at" = CURRENT_TIMESTAMP, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $1
	`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *userRepository) Recover(ctx context.Context, id int64) error {
	query := `
		UPDATE "users" 
		SET "deleted_at" = NULL, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $1
	`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *userRepository) CheckExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "users" WHERE "email" = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
