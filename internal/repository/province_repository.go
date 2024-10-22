package repository

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type provinceRepository struct {
	db database
}

func NewProvinceRepository(db database) ProvinceRepository {
	return &provinceRepository{
		db: db,
	}
}

func (r *provinceRepository) FindAll(ctx context.Context, name string) ([]model.Province, error) {
	var provinces []model.Province

	query := `SELECT "id","name" FROM "provinces" WHERE "name" ILIKE $1`
	rows, err := r.db.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var province model.Province
		if err := rows.Scan(&province.ID, &province.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}

		provinces = append(provinces, province)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return provinces, nil
}
