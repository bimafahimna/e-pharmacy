package repository

import (
	"context"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type manufacturerRepository struct {
	db database
}

func NewManufacturerRepository(db database) ManufacturerRepository {
	return &manufacturerRepository{
		db: db,
	}
}

func (r *manufacturerRepository) FindAll(ctx context.Context, name string) ([]model.Manufacturer, error) {
	query := `SELECT * FROM "manufacturers" WHERE "name" ILIKE $1`

	rows, err := r.db.QueryContext(ctx, query, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	manufacturers := []model.Manufacturer{}
	for rows.Next() {
		var manufacturer model.Manufacturer
		if err := rows.Scan(&manufacturer.ID, &manufacturer.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		manufacturers = append(manufacturers, manufacturer)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return manufacturers, nil
}

func (r *manufacturerRepository) CheckExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "manufacturers" WHERE "id" = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
