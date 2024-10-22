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

type pharmacistRepository struct {
	db database
}

func NewPharmacistRepository(db database) PharmacistRepository {
	return &pharmacistRepository{
		db: db,
	}
}

func (r *pharmacistRepository) Create(ctx context.Context, pharmacist model.Pharmacist) error {
	query := `
		INSERT INTO "pharmacist_details" (
			"user_id", "name", "sipa_number", "whatsapp_number", 
			"years_of_experience", "is_assigned"
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`
	args := []interface{}{
		pharmacist.UserID, pharmacist.Name, pharmacist.SipaNumber, pharmacist.WhatsappNumber,
		pharmacist.YearsOfExperience, pharmacist.IsAssigned,
	}

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

func (r *pharmacistRepository) FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Pharmacist, error) {
	query := `
		SELECT 
			p."user_id", p."name", u."email",
			p."sipa_number", p."whatsapp_number", p."years_of_experience", p."created_at", p."is_assigned"
		FROM "pharmacist_details" p
		JOIN "users" u ON u."id" = p."user_id"
		WHERE 1 = 1 
	`
	args := []interface{}{}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "email", "sipa_number", "whatsapp_number":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "years_of_experience", "is_assigned":
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

	pharmacists := []model.Pharmacist{}
	for rows.Next() {
		var pharmacist model.Pharmacist
		dest := []interface{}{
			&pharmacist.UserID, &pharmacist.Name, &pharmacist.Email,
			&pharmacist.SipaNumber, &pharmacist.WhatsappNumber, &pharmacist.YearsOfExperience, &pharmacist.CreatedAt, &pharmacist.IsAssigned,
		}

		if err := rows.Scan(dest...); err != nil {

			return nil, apperror.ErrInternalServerError
		}
		pharmacists = append(pharmacists, pharmacist)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return pharmacists, nil
}

func (r *pharmacistRepository) Find(ctx context.Context, userID int64) (*model.Pharmacist, error) {
	query := `
		SELECT 
			p."user_id", p."name", u."email", p."sipa_number",
			p."whatsapp_number", p."years_of_experience", p."is_assigned"
		FROM "pharmacist_details" p
		JOIN "users" u ON u."id" = p."user_id"
		WHERE "user_id" = $1
	`
	args := []interface{}{userID}

	var pharmacist model.Pharmacist
	dest := []interface{}{
		&pharmacist.UserID, &pharmacist.Name, &pharmacist.Email, &pharmacist.SipaNumber,
		&pharmacist.WhatsappNumber, &pharmacist.YearsOfExperience, &pharmacist.IsAssigned,
	}

	if err := r.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pharmacist, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &pharmacist, nil
}

func (r *pharmacistRepository) CountAll(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `
		SELECT COUNT(p.*)
		FROM "pharmacist_details" p
		JOIN "users" u ON u."id" = p."user_id"
		WHERE 1 = 1 
	`
	args := []interface{}{}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "email", "sipa_number", "whatsapp_number":
			query += fmt.Sprintf(` AND "%s" ILIKE $%d`, filter, len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "years_of_experience", "is_assigned":
			query += fmt.Sprintf(` AND "%s" = $%d`, filter, len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return total, nil
}

func (r *pharmacistRepository) Update(ctx context.Context, userID int64, pharmacist model.Pharmacist) error {
	query := `
		UPDATE "pharmacist_details"
		SET "name" = $1,
			"sipa_number" = $2,
			"whatsapp_number" = $3, 
			"years_of_experience" = $4,
			"is_assigned" = $5,
			"updated_at" = CURRENT_TIMESTAMP
		WHERE "user_id" = $6
	`
	args := []interface{}{
		&pharmacist.Name, &pharmacist.SipaNumber, &pharmacist.WhatsappNumber,
		&pharmacist.YearsOfExperience, &pharmacist.IsAssigned, userID,
	}

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

func (r *pharmacistRepository) Delete(ctx context.Context, userID int64) error {
	query := `
		DELETE FROM "pharmacist_details" WHERE "user_id" = $1
	`
	args := []interface{}{userID}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *pharmacistRepository) CheckAssigned(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT "is_assigned" FROM "pharmacist_details" WHERE "user_id" = $1`
	args := []interface{}{userID}

	var assigned bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&assigned); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperror.ErrNotFound
		}
		return false, apperror.ErrInternalServerError
	}
	return assigned, nil
}
