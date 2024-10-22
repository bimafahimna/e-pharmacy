package repository

import (
	"context"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type cityRepository struct {
	db database
}

func NewCityRepository(db database) CityRepository {
	return &cityRepository{
		db: db,
	}
}

func (r *cityRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]model.City, error) {
	var cities []model.City
	query := `
			SELECT 
				"id",
				"province_id",
				"province_unofficial_id",
				"name",
				"type",
				"unofficial_id" 
			FROM "cities" WHERE 1 = 1`

	args := []interface{}{}
	for filter, value := range filters {
		if value == nil {
			continue
		}
		switch filter {
		case "name":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "province_unofficial_id":
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
		var city model.City
		if err := rows.Scan(&city.ID, &city.ProvinceId, &city.ProvinceUnofficialId, &city.Name, &city.Type, &city.UnofficialId); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		cities = append(cities, city)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return cities, nil
}
