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

type productRepository struct {
	db database
}

func NewProductRepository(db database) ProductRepository {
	return &productRepository{
		db: db,
	}
}
func (r *productRepository) FindAllByPharmacyID(ctx context.Context, pharmacyID int, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Product, error) {
	query := `select 
			p."id",
			p."name",
			p."generic_name",
			m."name",
			pc."name",
			pf."name",
			p."is_active"
		from
			"products" p
		join "manufacturers" m on
			m."id" = p."manufacturer_id"
		join "product_classifications" pc on
			pc."id" = p."product_classification_id"
		join "product_forms" pf on
			pf."id" = p."product_form_id"
		where
			"deleted_at" is null
			and not exists (
			select
				*
			from
				pharmacy_products pp
			where
				p.id = pp.product_id
				and pp.pharmacy_id = $1)`

	aliases := map[string]string{
		"name":                   `p."name"`,
		"generic_name":           `p."generic_name"`,
		"manufacturer":           `m."name"`,
		"product_classification": `pc."name"`,
		"product_form":           `pf."name"`,
		"is_active":              `p."is_active"`,
		"created_at":             `p."created_at"`,
		"usage":                  `p."usage"`,
	}
	args := []interface{}{}
	args = append(args, pharmacyID)

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "generic_name", "manufacturer", "product_classification", "product_form":
			query += fmt.Sprintf(` AND %s ILIKE $%d`, aliases[filter], len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_active", "usage":
			query += fmt.Sprintf(` AND %s = $%d`, aliases[filter], len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	query += fmt.Sprintf(`
		ORDER BY %s %s
		OFFSET $%d
		LIMIT $%d
	`, aliases[sortBy], sort, len(args)+1, len(args)+2)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var product model.Product
		dest := []interface{}{
			&product.ID, &product.Name, &product.GenericName, &product.Manufacturer,
			&product.ProductClassification, &product.ProductForm, &product.IsActive,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return products, nil
}

func (r *productRepository) CountAllByPharmacyID(ctx context.Context, pharmacyID int, filters map[string]interface{}) (int, error) {
	query := `select 
			count(p.*)
		from
			"products" p
		join "manufacturers" m on
			m."id" = p."manufacturer_id"
		join "product_classifications" pc on
			pc."id" = p."product_classification_id"
		join "product_forms" pf on
			pf."id" = p."product_form_id"
		where
			"deleted_at" is null
			and not exists (
			select
				*
			from
				pharmacy_products pp
			where
				p.id = pp.product_id
				and pp.pharmacy_id = $1)`

	aliases := map[string]string{
		"name":                   `p."name"`,
		"generic_name":           `p."generic_name"`,
		"manufacturer":           `m."name"`,
		"product_classification": `pc."name"`,
		"product_form":           `pf."name"`,
		"is_active":              `p."is_active"`,
		"created_at":             `p."created_at"`,
		"usage":                  `p."usage"`,
	}
	args := []interface{}{}
	args = append(args, pharmacyID)

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "generic_name", "manufacturer", "product_classification", "product_form":
			query += fmt.Sprintf(` AND %s ILIKE $%d`, aliases[filter], len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_active", "usage":
			query += fmt.Sprintf(` AND %s = $%d`, aliases[filter], len(args)+1)
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

func (r *productRepository) InsertItem(ctx context.Context, product model.Product) error {
	query := `
		INSERT INTO "products" (
			"name",
			"manufacturer_id",
			"product_classification_id",
			"product_form_id",
			"generic_name",
			"categories",
			"description",
			"unit_in_pack",
			"selling_unit",
			"weight",
			"height",
			"length",
			"width",
			"image_url",
			"is_active",
			"usage",
			"product_ts"
		) 
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	if _, err := r.db.ExecContext(ctx, query,
		&product.Name,
		&product.ManufacturerID,
		&product.ProductClassificationID,
		&product.ProductFormID,
		&product.GenericName,
		&product.Categories,
		&product.Description,
		&product.UnitInPack,
		&product.SellingUnit,
		&product.Weight,
		&product.Height,
		&product.Length,
		&product.Width,
		&product.ImageURL,
		&product.IsActive,
		product.Usage,
		product.Agg,
	); err != nil {
		if postgres.IsForeignKeyViolation(err) {
			return apperror.ErrBadRequest
		}
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *productRepository) CheckExists(ctx context.Context, name string, genericName string, manufacturerId int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM "products"
			WHERE "name" = $1 AND "generic_name" = $2 AND "manufacturer_id" = $3
		)
	`

	var ok bool
	if err := r.db.QueryRowContext(ctx, query, name, genericName, manufacturerId).Scan(&ok); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, apperror.ErrInternalServerError
	}
	return ok, nil
}

func (r *productRepository) FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Product, error) {
	query := `
		SELECT 
			p."id", p."name", p."generic_name",
			m."name", pc."name", pf."name", p."is_active"
		FROM "products" p
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		WHERE "deleted_at" IS NULL
	`
	aliases := map[string]string{
		"name":                   `p."name"`,
		"generic_name":           `p."generic_name"`,
		"manufacturer":           `m."name"`,
		"product_classification": `pc."name"`,
		"product_form":           `pf."name"`,
		"is_active":              `p."is_active"`,
		"created_at":             `p."created_at"`,
		"usage":                  `p."usage"`,
	}
	args := []interface{}{}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "generic_name", "manufacturer", "product_classification", "product_form":
			query += fmt.Sprintf(` AND %s ILIKE $%d`, aliases[filter], len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_active", "usage":
			query += fmt.Sprintf(` AND %s = $%d`, aliases[filter], len(args)+1)
			args = append(args, value)
		default:
			continue
		}
	}

	query += fmt.Sprintf(`
		ORDER BY %s %s
		OFFSET $%d
		LIMIT $%d
	`, aliases[sortBy], sort, len(args)+1, len(args)+2)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var product model.Product
		dest := []interface{}{
			&product.ID, &product.Name, &product.GenericName, &product.Manufacturer,
			&product.ProductClassification, &product.ProductForm, &product.IsActive,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return products, nil
}

func (r *productRepository) CountAll(ctx context.Context, filters map[string]interface{}) (int, error) {
	query := `
		SELECT COUNT(*) FROM "products" p
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		WHERE "deleted_at" IS NULL
	`
	aliases := map[string]string{
		"name":                   `p."name"`,
		"generic_name":           `p."generic_name"`,
		"manufacturer":           `m."name"`,
		"product_classification": `pc."name"`,
		"product_form":           `pf."name"`,
		"is_active":              `p."is_active"`,
		"created_at":             `p."created_at"`,
		"usage":                  `p."usage"`,
	}
	args := []interface{}{}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "generic_name", "manufacturer", "product_classification", "product_form":
			query += fmt.Sprintf(` AND %s ILIKE $%d`, aliases[filter], len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_active", "usage":
			query += fmt.Sprintf(` AND %s = $%d`, aliases[filter], len(args)+1)
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

func (r *productRepository) IncrementUsage(ctx context.Context, id, diff int) error {
	query := `
		UPDATE "products"
		SET "usage" = "usage" + $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "id" = $2
	`
	args := []interface{}{diff, id}

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

func (r *productRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	var product model.Product
	query := `
	SELECT 
		p."name" ,
		p."generic_name" ,
		p."image_url" ,
		pf."name" as form, 
		m."name" as manufacturer,
		pc."name" as classification,
		p."categories" as categories ,
		p."description" ,
		p."selling_unit" ,
		p."unit_in_pack" ,
		p."weight" 
	FROM "products" p
	LEFT JOIN "manufacturers" m on m."id" = p."manufacturer_id"
	LEFT JOIN "product_classifications" pc on pc."id" = p."product_classification_id" 
	LEFT JOIN "product_forms" pf on pf."id" = p."product_form_id"
	WHERE p."id" = $1
	;
	`
	args := []interface{}{id}
	dest := []interface{}{
		&product.Name, &product.GenericName, &product.ImageURL,
		&product.ProductForm, &product.Manufacturer, &product.ProductClassification, &product.Categories,
		&product.Description, &product.SellingUnit, &product.UnitInPack, &product.Weight,
	}
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &product, nil
}

func (r *productRepository) FindIDBestseller(ctx context.Context, limit int) ([]int, error) {
	query := `
	SELECT "id"
	FROM "products"
	WHERE "is_active" = 'true'
	ORDER BY "daily_sold_amount" DESC,"sold_amount" DESC
	LIMIT $1
	`
	args := []interface{}{limit}
	return r.FindID(ctx, query, args)
}

func (r *productRepository) SearchID(ctx context.Context, search string, offset, limit int) ([]int, error) {
	query := `
	SELECT
		"id"
	FROM
		products
	WHERE
		product_ts @@ to_tsquery($1) AND "is_active" = 'true'
	ORDER BY "sold_amount" DESC
	OFFSET $2
	LIMIT $3
	`
	args := []interface{}{search, offset, limit}
	return r.FindID(ctx, query, args)
}

func (r *productRepository) FindID(ctx context.Context, query string, args []interface{}) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}

	var productId []int
	defer rows.Close()
	for rows.Next() {
		var id int
		dest := []interface{}{&id}
		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		productId = append(productId, id)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}

	return productId, nil
}
