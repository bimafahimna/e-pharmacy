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

type partnerRepository struct {
	db database
}

func NewPartnerRepository(db database) PartnerRepository {
	return &partnerRepository{
		db: db,
	}
}

func (r *partnerRepository) Create(ctx context.Context, partner model.Partner) error {
	query := `
	INSERT INTO "partners" 
		("name", "logo_url", "year_founded", "active_days", "operational_start","operational_stop","is_active") 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7)`

	if _, err := r.db.ExecContext(ctx, query, partner.Name, partner.LogoUrl, partner.YearFounded, partner.ActiveDays, partner.OperationalStart, partner.OperationalStop, partner.IsActive); err != nil {
		if postgres.IsUniqueViolation(err) {
			return apperror.ErrPartnerNameHasBeenRegistered
		}
		return apperror.ErrInternalServerError
	}
	return nil
}
func (r *partnerRepository) CountAll(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `
		SELECT
			COUNT(*)
		FROM "partners" p
		WHERE "deleted_at" IS NULL
	`
	args := []interface{}{}
	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "active_days":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "id", "year_founded", "operational_start", "operational_stop", "is_active":
			query += fmt.Sprintf(` AND "%s" = $%d`, filter, len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}
	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return total, apperror.ErrInternalServerError
	}

	return total, nil
}

func (r *partnerRepository) FindByName(ctx context.Context, name string) (*model.Partner, error) {
	var partner model.Partner
	query := `
		SELECT "id", "name", "logo_url", "year_founded", "active_days", "operational_start","operational_stop","is_active", "created_at", "updated_at" 
		FROM "partners" 
		WHERE "name" ILIKE $1 AND "deleted_at" IS NULL`

	args := []interface{}{name + "%"}
	dest := []interface{}{&partner.ID, &partner.Name, &partner.LogoUrl,
		&partner.YearFounded, &partner.ActiveDays, &partner.OperationalStart, &partner.OperationalStop, &partner.IsActive, &partner.CreatedAt, &partner.UpdatedAt}

	if err := r.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &partner, nil
}
func (r *partnerRepository) FindByID(ctx context.Context, id int) (*model.Partner, error) {
	var partner model.Partner
	query := `
		SELECT "id", "name", "logo_url", "year_founded", "active_days", "operational_start","operational_stop","is_active", "created_at", "updated_at" 
		FROM "partners" 
		WHERE "id" = $1 AND "deleted_at" IS NULL`

	args := []interface{}{id}
	dest := []interface{}{&partner.ID, &partner.Name, &partner.LogoUrl,
		&partner.YearFounded, &partner.ActiveDays, &partner.OperationalStart, &partner.OperationalStop, &partner.IsActive, &partner.CreatedAt, &partner.UpdatedAt}

	if err := r.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &partner, nil
}

func (r *partnerRepository) UpdateProfile(ctx context.Context, partner model.Partner) error {
	query := `
		UPDATE "partners"
		SET "name" = $1, "logo_url" = $2, "year_founded" = $3, "is_active" = $4, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $5 AND "deleted_at" IS NULL
	`
	args := []interface{}{partner.Name, partner.LogoUrl, partner.YearFounded, partner.IsActive, partner.ID}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return apperror.ErrUniqueViolation
		}
		return apperror.ErrInternalServerError
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *partnerRepository) UpdateDaysAndHours(ctx context.Context, partner model.Partner) error {
	query := `
		UPDATE "partners"
		SET "active_days" = $1, "operational_start" = $2, "operational_stop" = $3, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $4 AND "deleted_at" IS NULL
	`
	args := []interface{}{partner.ActiveDays, partner.OperationalStart, partner.OperationalStop, partner.ID}

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

func (r *partnerRepository) FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Partner, error) {
	query := `
		SELECT
			"id",
			"name",
			"logo_url",
			"year_founded",
			"active_days",
			"operational_start",
			"operational_stop",
			"is_active",
			"created_at"
		FROM "partners" p
		WHERE "deleted_at" IS NULL
	`
	args := []interface{}{}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "active_days":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "id", "year_founded", "operational_start", "operational_stop", "is_active":
			query += fmt.Sprintf(` AND "%s" = $%d`, filter, len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	query += fmt.Sprintf(`
		ORDER BY p."%s" %s
		OFFSET $%d
		LIMIT $%d
	`, sortBy, sort, len(args)+1, len(args)+2)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()
	partners := []model.Partner{}

	for rows.Next() {
		var p model.Partner
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.LogoUrl,
			&p.YearFounded,
			&p.ActiveDays,
			&p.OperationalStart,
			&p.OperationalStop,
			&p.IsActive,
			&p.CreatedAt,
		); err != nil {
			return nil, apperror.ErrInternalServerError
		}

		partners = append(partners, p)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return partners, nil
}
