package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository/postgres"
)

type pharmacyRepository struct {
	db database
}

func NewPharmacyRepository(db database) PharmacyRepository {
	return &pharmacyRepository{
		db: db,
	}
}

func (r *pharmacyRepository) Create(ctx context.Context, pharmacy model.Pharmacy) (int, error) {
	query := `
	INSERT INTO "pharmacies" 
		("pharmacist_id", "partner_id", "name", "address", "city_id","latitude","longitude","is_active") 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "id"
	`
	var id int
	if err := r.db.QueryRowContext(ctx, query, pharmacy.PharmacistId, pharmacy.PartnerId, pharmacy.Name, pharmacy.Address, pharmacy.CityId, pharmacy.Latitude, pharmacy.Longitude, pharmacy.IsActive).Scan(&id); err != nil {
		if postgres.IsUniqueViolation(err) {
			return 0, apperror.ErrPharmacyLocationAlreadyExist
		}
		return 0, apperror.ErrInternalServerError
	}
	return id, nil
}

func (r *pharmacyRepository) FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Pharmacy, error) {
	query := `
		SELECT 
			p1."id",
			p1."pharmacist_id",
			p3."name",
			p1."partner_id",
			p2."name",
			p1."name",
			p1."address",
			p1."city_id",
			c."name",
			p1."latitude",
			p1."longitude",
			p1."is_active",
			p1."created_at"
		FROM "pharmacies" p1
		JOIN "partners" p2 ON p2."id" = p1."partner_id"
		LEFT JOIN "users" u ON u."id" = p1."pharmacist_id"
		LEFT JOIN "pharmacist_details" p3 ON p3."user_id" = u."id"
		LEFT JOIN "cities" c ON c."unofficial_id" = p1."city_id"
		WHERE p1."deleted_at" IS NULL
	`
	args := []interface{}{}
	for filter, value := range filters {
		if value != "" {
			switch filter {
			case "name", "address":
				query += fmt.Sprintf(` AND p1."%s" ILIKE $%d`, filter, len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "pharmacist_name":
				query += fmt.Sprintf(" AND p3.\"name\" ILIKE $%d", len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "partner_name":
				query += fmt.Sprintf(" AND p2.\"name\" ILIKE $%d", len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "city_name":
				query += fmt.Sprintf(" AND c.\"name\" ILIKE $%d", len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "is_active":
				query += fmt.Sprintf(" AND p1.\"%s\" = $%d", filter, len(args)+1)
				args = append(args, value)
			default:
				continue
			}
		}
	}

	query += fmt.Sprintf(`
		ORDER BY p1."%s" %s
		OFFSET $%d
		LIMIT $%d
	`, sortBy, sort, len(args)+1, len(args)+2)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	pharmacies := []model.Pharmacy{}
	for rows.Next() {
		var p model.Pharmacy
		dest := []interface{}{
			&p.ID, &p.PharmacistId, &p.PharmacistName, &p.PartnerId, &p.PartnerName,
			&p.Name, &p.Address, &p.CityId, &p.CityName, &p.Latitude, &p.Longitude, &p.IsActive, &p.CreatedAt,
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		pharmacies = append(pharmacies, p)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return pharmacies, nil
}

func (r *pharmacyRepository) Find(ctx context.Context, pharmacistID int64) (*model.Pharmacy, error) {
	query := `
		SELECT 
			"id", "pharmacist_id", "partner_id", "name", "address", 
			"city_id", "latitude", "longitude", "is_active"
		FROM "pharmacies" WHERE "pharmacist_id" = $1
	`

	var pharmacy model.Pharmacy
	dest := []interface{}{
		&pharmacy.ID, &pharmacy.PharmacistId, &pharmacy.PartnerId, &pharmacy.Name, &pharmacy.Address,
		&pharmacy.CityId, &pharmacy.Latitude, &pharmacy.Longitude, &pharmacy.IsActive,
	}

	if err := r.db.QueryRowContext(ctx, query, pharmacistID).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &pharmacy, nil
}

func (r *pharmacyRepository) FindByID(ctx context.Context, pharmacyID int64) (*model.Pharmacy, error) {
	query := `
		SELECT 
			p."id", p."pharmacist_id", p."partner_id", p."name", p."address", 
			p."city_id", c."name", p."latitude", p."longitude", p."is_active"
		FROM "pharmacies" p
		LEFT JOIN "cities" c ON c."unofficial_id" = p."city_id"
		WHERE p."id" = $1
	`

	var pharmacy model.Pharmacy
	dest := []interface{}{
		&pharmacy.ID, &pharmacy.PharmacistId, &pharmacy.PartnerId, &pharmacy.Name, &pharmacy.Address,
		&pharmacy.CityId, &pharmacy.CityName, &pharmacy.Latitude, &pharmacy.Longitude, &pharmacy.IsActive,
	}

	if err := r.db.QueryRowContext(ctx, query, pharmacyID).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &pharmacy, nil
}

func (r *pharmacyRepository) CountAll(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `
	SELECT 
		COUNT(p1.*)
	FROM "pharmacies" p1
	JOIN "partners" p2 ON p2."id" = p1."partner_id"
	LEFT JOIN "users" u ON u."id" = p1."pharmacist_id"
	LEFT JOIN "pharmacist_details" p3 ON p3."user_id" = u."id"
	LEFT JOIN "cities" c ON c."unofficial_id" = p1."city_id"
	WHERE p1."deleted_at" IS NULL
`
	args := []interface{}{}
	for filter, value := range filters {
		if value != "" {
			switch filter {
			case "name", "address":
				query += fmt.Sprintf(` AND p1."%s" ILIKE $%d`, filter, len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "pharmacist_name":
				query += fmt.Sprintf(" AND p3.\"name\" ILIKE $%d", len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "partner_name":
				query += fmt.Sprintf(" AND p2.\"name\" ILIKE $%d", len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "city_name":
				query += fmt.Sprintf(" AND c.\"name\" ILIKE $%d", len(args)+1)
				args = append(args, fmt.Sprintf("%%%v%%", value))
			case "is_active":
				query += fmt.Sprintf(" AND p1.\"%s\" = $%d", filter, len(args)+1)
				args = append(args, value)
			default:
				continue
			}
		}
	}

	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return total, nil
}
