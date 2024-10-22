package repository

import (
	"context"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type productClassificationRepository struct {
	db database
}

func NewProductClassificationRepository(db database) ProductClassificationRepository {
	return &productClassificationRepository{
		db: db,
	}
}

func (r *productClassificationRepository) FindAll(ctx context.Context, name string) ([]model.ProductClassification, error) {
	query := `SELECT * FROM "product_classifications" WHERE "name" ILIKE $1`

	rows, err := r.db.QueryContext(ctx, query, fmt.Sprintf("%%%s%%", name))
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	classifications := []model.ProductClassification{}
	for rows.Next() {
		var classification model.ProductClassification
		if err := rows.Scan(&classification.ID, &classification.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		classifications = append(classifications, classification)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return classifications, nil
}

func (r *productClassificationRepository) CheckExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "product_classifications" WHERE "id" = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
