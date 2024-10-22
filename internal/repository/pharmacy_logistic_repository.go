package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type pharmacyLogisticRepository struct {
	db database
}

func NewPharmacyLogisticRepository(db database) PharmacyLogisticRepository {
	return &pharmacyLogisticRepository{
		db: db,
	}
}

func (r *pharmacyLogisticRepository) FindAllByPharmacyID(ctx context.Context, addressID, pharmacyID int) ([]model.PharmacyLogistic, error) {
	pharmaLogistics := []model.PharmacyLogistic{}

	query := `
	SELECT 
		l."id",
		l."name",
		l."logo_url",
		l."service",
		l."price_per_kilometers",
		l."eda",
		haversine_distance(ca."latitude",ca."longitude",p."latitude",p."longitude") as distance_km,
		ca."city_id",
		p."city_id",
		p."name"
	FROM "pharmacy_logistics" pl
	JOIN "logistics" l ON pl."logistic_id" = l."id"
	JOIN "pharmacies" p on p."id" = pl."pharmacy_id"
	JOIN "customer_addresses" ca ON ca."id" = $1
	WHERE pl."pharmacy_id"=$2
	`
	args := []interface{}{addressID, pharmacyID}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var pharmaLogistic model.PharmacyLogistic
		dest := []interface{}{
			&pharmaLogistic.Logistics.ID, &pharmaLogistic.Logistics.Name, &pharmaLogistic.Logistics.LogoUrl,
			&pharmaLogistic.Logistics.Service, &pharmaLogistic.Logistics.PricePerKilometers, &pharmaLogistic.Logistics.EDA,
			&pharmaLogistic.DistanceKM, &pharmaLogistic.CustomerCityID, &pharmaLogistic.PharmacyCityID, &pharmaLogistic.Pharmacies.Name,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		pharmaLogistics = append(pharmaLogistics, pharmaLogistic)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return pharmaLogistics, nil
}

func (r *pharmacyLogisticRepository) CheckEDA(ctx context.Context, logisticID, pharmacyID int) (int, error) {
	query := `
		SELECT l."eda" FROM "logistics" l
		JOIN "pharmacy_logistics" pl 
			ON pl."logistic_id" = $1 AND pl."pharmacy_id" = $2
	`
	args := []interface{}{logisticID, pharmacyID}

	var eda int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&eda); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperror.ErrNotFound
		}
		return 0, apperror.ErrInternalServerError
	}
	return eda, nil
}

func (r *pharmacyLogisticRepository) Create(ctx context.Context, logisticID, pharmacyID int) error {
	query := `
		INSERT INTO "pharmacy_logistics" (
			"pharmacy_id", "logistic_id"
		) VALUES ($1, $2) 
	`
	args := []interface{}{
		pharmacyID, logisticID,
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}
