package repository

import (
	"context"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type subDistrictRepository struct {
	db database
}

func NewSubDistrictRepository(db database) SubDistrictRepository {
	return &subDistrictRepository{
		db: db,
	}
}

func (r *subDistrictRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]model.SubDistrict, error) {
	var subDistricts []model.SubDistrict
	query := `SELECT "id", "district_id", "name" FROM "sub_districts" WHERE 1 = 1`

	args := []interface{}{}
	for filter, value := range filters {
		if value == nil {
			continue
		}
		switch filter {
		case "name":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "district_id":
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
		var subDistrict model.SubDistrict
		if err := rows.Scan(&subDistrict.ID, &subDistrict.DistrictId, &subDistrict.Name); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		subDistricts = append(subDistricts, subDistrict)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return subDistricts, nil
}
