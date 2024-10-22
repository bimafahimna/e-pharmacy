package repository

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type productCategoryRepository struct {
	db database
}

func NewProductCategoryRepository(db database) ProductCategoryRepository {
	return &productCategoryRepository{
		db: db,
	}
}

func (r *productCategoryRepository) FindAll(ctx context.Context) ([]model.ProductCategory, error) {
	query := `SELECT "id", "name" FROM "categories"`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}

	categories := []model.ProductCategory{}
	for rows.Next() {
		var category model.ProductCategory
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}

		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return categories, nil
}

func (r *productCategoryRepository) Create(ctx context.Context, category model.ProductCategory) error {
	query := `
		INSERT INTO "categories" ("name") VALUES ($1)
		ON CONFLICT DO NOTHING
	`
	args := []interface{}{category.Name}

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

func (r *productCategoryRepository) Update(ctx context.Context, id int, category model.ProductCategory) error {
	query := `
		UPDATE "categories"
		SET "name" = $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $2
	`
	args := []interface{}{category.Name, id}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *productCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM "categories" WHERE "id" = $1`
	args := []interface{}{id}

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

func (r *productCategoryRepository) CheckExists(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM "categories" WHERE "name" = $1)`
	args := []interface{}{name}

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}
