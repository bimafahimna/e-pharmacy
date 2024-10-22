package repository

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type logisticRepository struct {
	db database
}

func NewLogisticRepository(db database) LogisticRepository {
	return &logisticRepository{db: db}
}

func (r *logisticRepository) FindAll(ctx context.Context) ([]model.Logistic, error) {
	query := `
	SELECT
		"id",
		"name",
		"logo_url", 
		"service", 
		"price_per_kilometers",
		"eda" 
	FROM "logistics"
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	var logistics []model.Logistic
	for rows.Next() {
		var logistic model.Logistic
		dest := []interface{}{
			&logistic.ID, &logistic.Name, &logistic.LogoUrl, &logistic.Service,
			&logistic.PricePerKilometers, &logistic.EDA,
		}

		if err := rows.Scan(dest...); err != nil {

			return nil, apperror.ErrInternalServerError
		}
		logistics = append(logistics, logistic)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return logistics, nil
}
