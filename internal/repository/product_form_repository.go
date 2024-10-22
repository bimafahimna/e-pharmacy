package repository

import (
	"context"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type productFormRepository struct {
	db database
}

func NewProductFormRepository(db database) ProductFormRepository {
	return &productFormRepository{
		db: db,
	}
}

func (r *productFormRepository) FindAll(ctx context.Context, name string) ([]model.ProductForm, error) {
	query := `SELECT * FROM "product_forms" WHERE "name" ILIKE $1`

	rows, err := r.db.QueryContext(ctx, query, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	forms := []model.ProductForm{}
	for rows.Next() {
		var form model.ProductForm
		if err := rows.Scan(&form.ID, &form.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		forms = append(forms, form)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return forms, nil
}

func (r *productFormRepository) CheckExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "product_forms" WHERE "id" = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
