package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type Store interface {
	Atomic(ctx context.Context, fn func(Store) error) error
	User() UserRepository
	VerificationToken() VerificationTokenRepository
	PasswordResetToken() PasswordResetTokenRepository
	SubDistrict() SubDistrictRepository
	District() DistrictRepository
	City() CityRepository
	Province() ProvinceRepository
	CustomerAddress() CustomerAddressRepository
	Partner() PartnerRepository
	Pharmacist() PharmacistRepository
	Pharmacy() PharmacyRepository
	ProductCategory() ProductCategoryRepository
	Manufacturer() ManufacturerRepository
	ProductClassification() ProductClassificationRepository
	ProductForm() ProductFormRepository
	Product() ProductRepository
	PharmacyProduct() PharmacyProductRepository
	Cart() CartRepository
	Order() OrderRepository
	OrderItem() OrderItemRepository
	Payment() PaymentRepository
	PharmacyLogistic() PharmacyLogisticRepository
	LogisticRepository() LogisticRepository
}

type store struct {
	conn *sql.DB
	db   database
}

func NewStore(db *sql.DB) Store {
	return &store{
		conn: db,
		db:   db,
	}
}

func (s *store) Atomic(ctx context.Context, fn func(Store) error) error {
	tx, err := s.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return apperror.ErrInternalServerError
	}

	defer commitOrRollback(tx, recover(), &err)

	err = fn(&store{conn: s.conn, db: tx})
	return err
}

func commitOrRollback(tx *sql.Tx, p interface{}, err *error) {
	if p != nil {
		tx.Rollback()
		panic(p)
	}

	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			*err = fmt.Errorf("tx err: %v, rollback err: %v", err, rollbackErr)
		}
		return
	}

	*err = tx.Commit()
}

func (s *store) User() UserRepository {
	return NewUserRepository(s.db)
}

func (s *store) VerificationToken() VerificationTokenRepository {
	return NewVerifyTokenRepository(s.db)
}

func (s *store) PasswordResetToken() PasswordResetTokenRepository {
	return NewPasswordResetTokenRepository(s.db)
}

func (s *store) SubDistrict() SubDistrictRepository {
	return NewSubDistrictRepository(s.db)
}

func (s *store) District() DistrictRepository {
	return NewDistrictRepository(s.db)
}

func (s *store) City() CityRepository {
	return NewCityRepository(s.db)
}

func (s *store) Province() ProvinceRepository {
	return NewProvinceRepository(s.db)
}

func (s *store) CustomerAddress() CustomerAddressRepository {
	return NewCustomerAddressRepository(s.db)
}

func (s *store) Partner() PartnerRepository {
	return NewPartnerRepository(s.db)
}

func (s *store) Pharmacist() PharmacistRepository {
	return NewPharmacistRepository(s.db)
}

func (s *store) Pharmacy() PharmacyRepository {
	return NewPharmacyRepository(s.db)
}

func (s *store) ProductCategory() ProductCategoryRepository {
	return NewProductCategoryRepository(s.db)
}

func (s *store) Manufacturer() ManufacturerRepository {
	return NewManufacturerRepository(s.db)
}

func (s *store) ProductClassification() ProductClassificationRepository {
	return NewProductClassificationRepository(s.db)
}

func (s *store) ProductForm() ProductFormRepository {
	return NewProductFormRepository(s.db)
}

func (s *store) Product() ProductRepository {
	return NewProductRepository(s.db)
}

func (s *store) PharmacyProduct() PharmacyProductRepository {
	return NewPharmacyProductRepository(s.db)
}

func (s *store) Cart() CartRepository {
	return NewCartRepository(s.db)
}

func (s *store) Order() OrderRepository {
	return NewOrderRepository(s.db)
}

func (s *store) OrderItem() OrderItemRepository {
	return NewOrderItemRepository(s.db)
}

func (s *store) Payment() PaymentRepository {
	return NewPaymentRepository(s.db)
}

func (s *store) PharmacyLogistic() PharmacyLogisticRepository {
	return NewPharmacyLogisticRepository(s.db)
}

func (s *store) LogisticRepository() LogisticRepository {
	return NewLogisticRepository(s.db)
}
