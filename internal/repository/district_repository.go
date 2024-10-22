package repository

import (
	"context"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type districtRepository struct {
	db database
}

func NewDistrictRepository(db database) DistrictRepository {
	return &districtRepository{
		db: db,
	}
}

func (r *districtRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]model.District, error) {
	var districts []model.District
	query := `SELECT "id", "city_id", "name" FROM "districts" WHERE 1 = 1`

	args := []interface{}{}
	for filter, value := range filters {
		if value == nil {
			continue
		}
		switch filter {
		case "name":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "city_id":
			query += fmt.Sprintf(` AND "%s" = $%d`, filter, len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var district model.District
		if err := rows.Scan(&district.ID, &district.CityId, &district.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		districts = append(districts, district)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return districts, nil
}
