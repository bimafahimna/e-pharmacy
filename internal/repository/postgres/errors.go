package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var pgErr *pgconn.PgError

var (
	CodeForeignKeyConstraintViolation = "23503"
	CodeUniqueConstraintViolation     = "23505"
)

func IsForeignKeyViolation(err error) bool {
	if errors.As(err, &pgErr); pgErr.Code == CodeForeignKeyConstraintViolation {
		return true
	}
	return false
}

func IsUniqueViolation(err error) bool {
	if errors.As(err, &pgErr); pgErr.Code == CodeUniqueConstraintViolation {
		return true
	}
	return false
}
