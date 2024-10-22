package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/shopspring/decimal"
)

type pharmacyProductRepository struct {
	db database
}

func NewPharmacyProductRepository(db database) PharmacyProductRepository {
	return &pharmacyProductRepository{
		db: db,
	}
}

func (r *pharmacyProductRepository) FindAllBestseller(ctx context.Context, bestSellerID []int) ([]model.Item, error) {
	query := `
	WITH "bestseller" AS (
		SELECT
			pp."pharmacy_id", pp."product_id", p."product_form_id",
			p."name", pp."price", pp."stock", p."selling_unit", p."image_url",
			ROW_NUMBER() OVER (PARTITION BY pp."product_id" ORDER BY l."eda" ASC,pp."sold_amount" DESC,pp."price" DESC) AS "rn"
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "pharmacies" ph ON ph."id" = pp."pharmacy_id"
		JOIN "partners" pa ON pa."id" = ph."partner_id"
		JOIN "pharmacy_logistics" pl ON pl."pharmacy_id" = ph."id"
		JOIN "logistics" l ON l."id" = pl."logistic_id"
		WHERE pp."stock" > 0 AND pp."price" != 0 AND ph."is_active" = 'true' AND pa."is_active" = 'true'
			AND pp."product_id" = ANY($1)
	)
	SELECT
		"pharmacy_id", "product_id", "product_form_id",
		"name", "price", "stock", "selling_unit", "image_url"
	FROM "bestseller" WHERE "rn" <= '1';
	`
	args := []interface{}{bestSellerID}
	return r.FindAll(ctx, query, args...)
}

func (r *pharmacyProductRepository) FindAllBestsellerPharmacies(ctx context.Context, productID, limit int) ([]model.AvailablePharmacy, error) {
	query := `
		SELECT
			pp."pharmacy_id", p."name" as "pharmacy_name",pt."name"as"partner_name",pt."logo_url"as"partner_logo",
			ph."address",c."name",pp."product_id", pp."price", pp."stock"
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "pharmacies" ph ON ph."id" = pp."pharmacy_id"
		JOIN "partners" pt ON pt."id" = ph."partner_id"
		JOIN "cities" c ON c."unofficial_id" = ph."city_id"
		WHERE pp."stock" > 0
			AND pp."product_id" = $1
		ORDER BY pp."sold_amount" DESC
		LIMIT $2
	`
	args := []interface{}{productID, limit}
	return r.FindAllPharmacies(ctx, query, args...)
}

func (r *pharmacyProductRepository) FindAllRecommended(ctx context.Context, productIDs []int, minLat, maxLat, minLong, maxLong float64) ([]model.Item, error) {
	query := `
		WITH "recommended" AS (
			SELECT
				pp."pharmacy_id", pp."product_id", p."product_form_id",
				p."name", pp."price", pp."stock", p."selling_unit", p."image_url",
				ROW_NUMBER() OVER (PARTITION BY pp."product_id" ORDER BY l."eda" ASC, pp."sold_amount" DESC) AS "rn"
			FROM "pharmacy_products" pp
			JOIN "products" p ON p."id" = pp."product_id"
			JOIN "pharmacies" ph ON ph."id" = pp."pharmacy_id"
			JOIN "partners" pa ON pa."id" = ph."partner_id"
			JOIN "pharmacy_logistics" pl ON pl."pharmacy_id" = ph."id"
			JOIN "logistics" l ON l."id" = pl."logistic_id"
			WHERE pp."stock" > 0 AND pp."price" != 0 AND ph."is_active" = 'true' AND pa."is_active" = 'true'
				AND pp."product_id" = ANY($1) AND ph."latitude" BETWEEN $2 AND $3 AND ph."longitude" BETWEEN $4 AND $5
		)
		SELECT
			"pharmacy_id", "product_id", "product_form_id",
			"name", "price", "stock", "selling_unit", "image_url"
		FROM "recommended" WHERE "rn" <= '1';
	`
	args := []interface{}{productIDs, minLat, maxLat, minLong, maxLong}
	return r.FindAll(ctx, query, args...)
}

func (r *pharmacyProductRepository) Search(ctx context.Context, offset, limit, maxDistance int, search string, latitude, longitude decimal.Decimal) ([]model.Item, error) {
	query := `
		SELECT 	
			pp."pharmacy_id", pp."product_id", p."product_form_id",
			p."name", pp."price", pp."stock", p."selling_unit", p."image_url"
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		JOIN "pharmacies" ph ON ph."id" = pp."pharmacy_id"
		WHERE pp."deleted_at" IS NULL AND pp."is_active" = true AND p."is_active" = true AND (
			m."name" ILIKE $1 OR pc."name" ILIKE $1 OR pf."name" ILIKE $1
			OR p."name" ILIKE $1 OR p."generic_name" ILIKE $1 OR p."description" ILIKE $1
			OR p."selling_unit" ILIKE $1 OR EXISTS (
				SELECT 1 FROM UNNEST(p."categories") AS "category" WHERE "category" ILIKE $1
			)
		) AND pp."stock" > 0
	`

	args := []interface{}{fmt.Sprintf("%%%s%%", search)}

	if !latitude.IsZero() && !longitude.IsZero() {
		distance := fmt.Sprintf(`HAVERSINE_DISTANCE($%d, $%d, ph."latitude", ph."longitude")`, len(args)+1, len(args)+2)
		args = append(args, latitude, longitude)

		query += fmt.Sprintf(`
			AND %s < $%d
			ORDER BY %s
		`, distance, len(args)+1, distance)
		args = append(args, maxDistance)
	}

	query += fmt.Sprintf(`
		OFFSET $%d
		LIMIT $%d
	`, len(args)+1, len(args)+2)
	args = append(args, offset, limit)
	return r.FindAll(ctx, query, args...)
}

func (r *pharmacyProductRepository) FindAll(ctx context.Context, query string, args ...interface{}) ([]model.Item, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	items := []model.Item{}
	for rows.Next() {
		var item model.Item
		dest := []interface{}{
			&item.PharmacyID, &item.ProductID, &item.ProductFormID,
			&item.Name, &item.Price, &item.Stock, &item.SellingUnit, &item.ImageURL,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return items, nil
}

func (r *pharmacyProductRepository) FindAllPharmacies(ctx context.Context, query string, args ...interface{}) ([]model.AvailablePharmacy, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	pharmacyProducts := []model.AvailablePharmacy{}
	for rows.Next() {
		var pharmacyProduct model.AvailablePharmacy
		dest := []interface{}{
			&pharmacyProduct.PharmacyID, &pharmacyProduct.PharmacyName, &pharmacyProduct.PartnerName, &pharmacyProduct.PartnerLogo,
			&pharmacyProduct.Address, &pharmacyProduct.CityName, &pharmacyProduct.ProductID, &pharmacyProduct.Price,
			&pharmacyProduct.Stock,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		pharmacyProducts = append(pharmacyProducts, pharmacyProduct)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return pharmacyProducts, nil
}

func (r *pharmacyProductRepository) FindAllByPharmacyID(ctx context.Context, pharmacyID int, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.PharmacyProduct, error) {
	query := `
		SELECT 
			pp."pharmacy_id", pp."product_id", pp."stock", pp."price", p."name",
			p."generic_name", m."name", pc."name", pf."name", pp."is_active"
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		WHERE pp."pharmacy_id" = $1 AND pp."deleted_at" IS NULL
	`
	aliases := map[string]string{
		"name":                   `p."name"`,
		"generic_name":           `p."generic_name"`,
		"manufacturer":           `m."name"`,
		"product_classification": `pc."name"`,
		"product_form":           `pf."name"`,
		"is_active":              `pp."is_active"`,
		"created_at":             `pp."created_at"`,
		"stock":                  `pp."stock"`,
	}
	args := []interface{}{pharmacyID}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "generic_name", "manufacturer", "product_classification", "product_form":
			query += fmt.Sprintf(` AND %s ILIKE $%d`, aliases[filter], len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_active":
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

	items := []model.PharmacyProduct{}
	for rows.Next() {
		var item model.PharmacyProduct
		dest := []interface{}{
			&item.PharmacyID, &item.ProductID, &item.Stock, &item.Price, &item.Name, &item.GenericName,
			&item.Manufacturer, &item.ProductClassification, &item.ProductForm, &item.IsActive,
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return items, nil
}

func (r *pharmacyProductRepository) Find(ctx context.Context, pharmacyID, productID int) (*model.PharmacyProduct, error) {
	query := `
		SELECT
			pp."pharmacy_id", pp."product_id", p."name", p."generic_name",
			m."name", pc."name", pf."name", pp."stock", pp."is_active", pp."price", pp."sold_amount"
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		WHERE pp."pharmacy_id" = $1 AND pp."product_id" = $2
	`
	args := []interface{}{pharmacyID, productID}

	var item model.PharmacyProduct
	dest := []interface{}{
		&item.PharmacyID, &item.ProductID, &item.Name, &item.GenericName,
		&item.Manufacturer, &item.ProductClassification, &item.ProductForm, &item.Stock,
		&item.IsActive, &item.Price, &item.SoldAmount,
	}

	if err := r.db.QueryRowContext(ctx, query, args...).Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternalServerError
	}
	return &item, nil
}

func (r *pharmacyProductRepository) CountSearch(ctx context.Context, maxDistance int, search string, latitude, longitude decimal.Decimal) (int, error) {
	query := `
		SELECT COUNT(pp.*)
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		JOIN "pharmacies" ph ON ph."id" = pp."pharmacy_id"
		WHERE pp."deleted_at" IS NULL AND pp."is_active" = true AND p."is_active" = true AND (
			m."name" ILIKE $1 OR pc."name" ILIKE $1 OR pf."name" ILIKE $1
			OR p."name" ILIKE $1 OR p."generic_name" ILIKE $1 OR p."description" ILIKE $1
			OR p."selling_unit" ILIKE $1 OR EXISTS (
				SELECT 1 FROM UNNEST(p."categories") AS "category" WHERE "category" ILIKE $1
			)
		) AND pp."stock" > 0
	`
	args := []interface{}{fmt.Sprintf("%%%s%%", search)}

	if !latitude.IsZero() && !longitude.IsZero() {
		query += fmt.Sprintf(`AND HAVERSINE_DISTANCE($%d, $%d, ph."latitude", ph."longitude") < $%d`, len(args)+1, len(args)+2, len(args)+3)
		args = append(args, latitude, longitude, maxDistance)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, apperror.ErrInternalServerError
	}
	return total, nil
}

func (r *pharmacyProductRepository) CountAllByPharmacyID(ctx context.Context, pharmacyID int, filters map[string]interface{}) (int, error) {
	query := `
		SELECT COUNT(pp.*)
		FROM "pharmacy_products" pp
		JOIN "products" p ON p."id" = pp."product_id"
		JOIN "manufacturers" m ON m."id" = p."manufacturer_id"
		JOIN "product_classifications" pc ON pc."id" = p."product_classification_id"
		JOIN "product_forms" pf ON pf."id" = p."product_form_id"
		WHERE pp."pharmacy_id" = $1 AND pp."deleted_at" IS NULL
	`
	aliases := map[string]string{
		"name":                   `p."name"`,
		"generic_name":           `p."generic_name"`,
		"manufacturer":           `m."name"`,
		"product_classification": `pc."name"`,
		"product_form":           `pf."name"`,
		"is_active":              `pp."is_active"`,
		"created_at":             `pp."created_at"`,
		"stock":                  `pp."stock"`,
	}
	args := []interface{}{pharmacyID}

	for filter, value := range filters {
		if value == nil {
			continue
		}

		switch filter {
		case "name", "generic_name", "manufacturer", "product_classification", "product_form":
			query += fmt.Sprintf(` AND %s ILIKE $%d`, aliases[filter], len(args)+1)
			args = append(args, fmt.Sprintf("%%%v%%", value))
		case "is_active":
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

func (r *pharmacyProductRepository) Insert(ctx context.Context, pharmacyID int, item model.PharmacyProduct) error {
	query := `
		INSERT INTO "pharmacy_products" (
			"pharmacy_id", "product_id", "sold_amount", "price", "stock", "is_active"
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	args := []interface{}{pharmacyID, item.ProductID, item.SoldAmount, item.Price, item.Stock, item.IsActive}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *pharmacyProductRepository) Update(ctx context.Context, pharmacyID, productID int, item model.PharmacyProduct) error {
	query := `
		UPDATE "pharmacy_products"
		SET "stock" = $1, "is_active" = $2, "updated_at" = CURRENT_TIMESTAMP
		WHERE "pharmacy_id" = $3 AND "product_id" = $4
	`
	args := []interface{}{item.Stock, item.IsActive, pharmacyID, productID}

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

func (r *pharmacyProductRepository) UpdateStock(ctx context.Context, pharmacyID, productID, diff int) error {
	query := `
		UPDATE "pharmacy_products"
		SET "stock" = "stock" + $1, "updated_at" = CURRENT_TIMESTAMP
		WHERE "pharmacy_id" = $2 AND "product_id" = $3
	`
	args := []interface{}{diff, pharmacyID, productID}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *pharmacyProductRepository) Delete(ctx context.Context, pharmacyID, productID int) error {
	query := `
		DELETE FROM "pharmacy_products" 
		WHERE "pharmacy_id" = $1 AND "product_id" = $2 AND "deleted_at" IS NULL
	`
	args := []interface{}{pharmacyID, productID}

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

func (r *pharmacyProductRepository) CheckExists(ctx context.Context, pharmacyID, productID int) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM "pharmacy_products"
		WHERE "pharmacy_id" = $1 AND "product_id" = $2
	)`
	args := []interface{}{pharmacyID, productID}

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, apperror.ErrInternalServerError
	}
	return exists, nil
}

func (r *pharmacyProductRepository) CheckUpdatedToday(ctx context.Context, pharmacyID, productID int) (bool, error) {
	query := `
		SELECT "updated_at"::DATE = CURRENT_TIMESTAMP::DATE FROM "pharmacy_products"
		WHERE "pharmacy_id" = $1 AND "product_id" = $2
	`
	args := []interface{}{pharmacyID, productID}

	var updatedToday bool
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&updatedToday); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperror.ErrNotFound
		}
		return false, apperror.ErrInternalServerError
	}
	return updatedToday, nil
}

func (r *pharmacyProductRepository) CheckSold(ctx context.Context, pharmacyID, productID int) (bool, error) {
	query := `
		SELECT "sold_amount" FROM "pharmacy_products"
		WHERE "pharmacy_id" = $1 AND "product_id" = $2
	`
	args := []interface{}{pharmacyID, productID}

	var soldAmount int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&soldAmount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, apperror.ErrNotFound
		}
		return false, apperror.ErrInternalServerError
	}

	if soldAmount == 0 {
		return false, nil
	}
	return true, nil
}

func (r *pharmacyProductRepository) FanOutFanInBestSeller(ctx context.Context, data []int, worker int) ([]model.Item, error) {
	memo := make([][]model.Item, worker)
	errors := make([]error, worker)
	deltaFloat := float64(len(data)) / float64(worker)
	delta := math.Ceil(deltaFloat)
	start := 0
	var wg sync.WaitGroup
	for i := 0; i < worker-1; i++ {
		end := start + int(delta)
		wg.Add(1)
		go func(index int, dataId []int, wg *sync.WaitGroup) {
			defer wg.Done()
			res, err := r.FindAllBestseller(ctx, dataId)
			memo[index] = res
			errors[index] = err
		}(i, data[start:end], &wg)
		start = end
	}
	wg.Add(1)
	go func(index int, dataId []int, wg *sync.WaitGroup) {
		defer wg.Done()
		res, err := r.FindAllBestseller(ctx, dataId)
		memo[index] = res
		errors[index] = err
	}(worker-1, data[start:], &wg)
	wg.Wait()
	var res []model.Item
	for _, data := range memo {
		res = append(res, data...)
	}
	if errors[0] != nil {
		return nil, apperror.ErrInternalServerError
	}
	return res, nil
}

func (r *pharmacyProductRepository) FanOutFanInRecommended(ctx context.Context, data []int, minLat, maxLat, minLong, maxLong float64, worker int) ([]model.Item, error) {
	memo := make([][]model.Item, worker)
	errors := make([]error, worker)
	deltaFloat := float64(len(data)) / float64(worker)
	delta := math.Ceil(deltaFloat)
	start := 0
	var wg sync.WaitGroup
	for i := 0; i < worker-1; i++ {
		end := start + int(delta)
		wg.Add(1)
		go func(index int, dataId []int, wg *sync.WaitGroup) {
			defer wg.Done()
			res, err := r.FindAllRecommended(ctx, dataId, minLat, maxLat, minLong, maxLong)
			memo[index] = res
			errors[index] = err
		}(i, data[start:end], &wg)
		start = end
	}
	wg.Add(1)
	go func(index int, dataId []int, wg *sync.WaitGroup) {
		defer wg.Done()
		res, err := r.FindAllRecommended(ctx, dataId, minLat, maxLat, minLong, maxLong)
		memo[index] = res
		errors[index] = err
	}(worker-1, data[start:], &wg)
	wg.Wait()
	var res []model.Item
	for _, data := range memo {
		res = append(res, data...)
	}
	if errors[0] != nil {
		return nil, apperror.ErrInternalServerError
	}
	return res, nil
}
