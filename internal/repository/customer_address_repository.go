package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository/postgres"
)

type customerAddressRepository struct {
	db database
}

func NewCustomerAddressRepository(db database) CustomerAddressRepository {
	return &customerAddressRepository{
		db: db,
	}
}

func (r *customerAddressRepository) Create(ctx context.Context, address model.CustomerAddress) error {
	query := `
		INSERT INTO "customer_addresses" 
			(
			"user_id",
			"name",
			"receiver_name",
			"receiver_phone_number",
			"latitude",
			"longitude",
			"province",
			"city",
			"district",
			"sub_district",
			"address_details",
			"is_active",
			"city_id"
			) 
		VALUES 
			(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
			)
	`
	if _, err := r.db.ExecContext(ctx, query, address.UserID, address.Name, address.ReceiverName, address.ReceiverPhoneNumber, address.Latitude, address.Longitude, address.Province, address.City, address.District, address.SubDistrict, address.AddressDetails, address.IsActive, address.CityID); err != nil {
		if postgres.IsUniqueViolation(err) {
			return apperror.ErrAddressNameAlreadyExist
		}
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *customerAddressRepository) FindActiveByUserId(ctx context.Context, userId int64) (*model.CustomerAddress, error) {
	var address model.CustomerAddress
	address.UserID = userId

	query := `
		SELECT "id", "name","receiver_name","receiver_phone_number", "latitude", "longitude", "province","city","district","sub_district","address_details","is_active","created_at","updated_at"
		FROM "customer_addresses"
		WHERE "is_active" = true AND "user_id" = $1;
	`
	if err := r.db.QueryRowContext(ctx, query, userId).Scan(&address.ID, &address.Name, &address.ReceiverName, &address.ReceiverPhoneNumber, &address.Latitude, &address.Longitude,
		&address.Province, &address.City, &address.District, &address.SubDistrict, &address.AddressDetails, &address.IsActive,
		&address.CreatedAt, &address.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperror.ErrInternalServerError
	}
	return &address, nil
}

func (r *customerAddressRepository) UnsetActive(ctx context.Context, id int64) error {
	query := `UPDATE "customer_addresses" SET "is_active" = false where "id" = $1;`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return apperror.ErrInternalServerError
	}
	return nil
}

func (r *customerAddressRepository) FindAllByUserId(ctx context.Context, userId int64) ([]model.CustomerAddress, error) {
	var addresses []model.CustomerAddress
	query := `
		SELECT "id", "name", "receiver_name","receiver_phone_number", "latitude", "longitude", "province","city","district","sub_district","address_details","is_active","created_at","updated_at"
		FROM "customer_addresses"
		WHERE "user_id" = $1;
	`
	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}
	defer rows.Close()

	for rows.Next() {
		var address model.CustomerAddress
		if err := rows.Scan(&address.ID, &address.Name, &address.ReceiverName, &address.ReceiverPhoneNumber, &address.Latitude, &address.Longitude,
			&address.Province, &address.City, &address.District, &address.SubDistrict, &address.AddressDetails, &address.IsActive,
			&address.CreatedAt, &address.UpdatedAt); err != nil {
			return nil, apperror.ErrInternalServerError
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.ErrInternalServerError
	}
	return addresses, nil
}
