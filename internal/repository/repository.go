package repository

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/shopspring/decimal"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (int64, error)
	FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.User, error)
	CountAll(ctx context.Context, filters map[string]interface{}) (int, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Verify(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, id int64, passwordHash string) error
	Delete(ctx context.Context, id int64) error
	Recover(ctx context.Context, id int64) error
	CheckExists(ctx context.Context, email string) (bool, error)
}

type VerificationTokenRepository interface {
	Save(ctx context.Context, userID int64, token string) error
	Find(ctx context.Context, token string) (int64, error)
	Revoke(ctx context.Context, userID int64) error
	CheckExists(ctx context.Context, userID int64) (bool, error)
}

type PasswordResetTokenRepository interface {
	Save(ctx context.Context, userID int64, token string) error
	FindUserID(ctx context.Context, token string) (int64, error)
	CheckExists(ctx context.Context, userID int64) (bool, error)
	Revoke(ctx context.Context, userID int64) error
}

type ProvinceRepository interface {
	FindAll(ctx context.Context, name string) ([]model.Province, error)
}

type CityRepository interface {
	FindAll(ctx context.Context, filters map[string]interface{}) ([]model.City, error)
}

type DistrictRepository interface {
	FindAll(ctx context.Context, filters map[string]interface{}) ([]model.District, error)
}

type SubDistrictRepository interface {
	FindAll(ctx context.Context, filters map[string]interface{}) ([]model.SubDistrict, error)
}

type CustomerAddressRepository interface {
	Create(ctx context.Context, address model.CustomerAddress) error
	FindActiveByUserId(ctx context.Context, userId int64) (*model.CustomerAddress, error)
	UnsetActive(ctx context.Context, id int64) error
	FindAllByUserId(ctx context.Context, userId int64) ([]model.CustomerAddress, error)
}

type PartnerRepository interface {
	Create(ctx context.Context, partner model.Partner) error
	FindByID(ctx context.Context, id int) (*model.Partner, error)
	FindByName(ctx context.Context, name string) (*model.Partner, error)
	FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Partner, error)
	CountAll(ctx context.Context, filter map[string]interface{}) (int, error)
	UpdateProfile(ctx context.Context, partner model.Partner) error
	UpdateDaysAndHours(ctx context.Context, partner model.Partner) error
}

type PharmacistRepository interface {
	Create(ctx context.Context, pharmacist model.Pharmacist) error
	FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Pharmacist, error)
	Find(ctx context.Context, userID int64) (*model.Pharmacist, error)
	CountAll(ctx context.Context, filters map[string]interface{}) (int, error)
	Update(ctx context.Context, userID int64, pharmacist model.Pharmacist) error
	Delete(ctx context.Context, userID int64) error
	CheckAssigned(ctx context.Context, userID int64) (bool, error)
}

type PharmacyRepository interface {
	Create(ctx context.Context, pharmacy model.Pharmacy) (int, error)
	FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Pharmacy, error)
	CountAll(ctx context.Context, filters map[string]interface{}) (int, error)
	Find(ctx context.Context, pharmacistID int64) (*model.Pharmacy, error)
	FindByID(ctx context.Context, pharmacyID int64) (*model.Pharmacy, error)
}

type ProductCategoryRepository interface {
	FindAll(ctx context.Context) ([]model.ProductCategory, error)
	Create(ctx context.Context, category model.ProductCategory) error
	Update(ctx context.Context, id int, category model.ProductCategory) error
	Delete(ctx context.Context, id int) error
	CheckExists(ctx context.Context, name string) (bool, error)
}

type ManufacturerRepository interface {
	FindAll(ctx context.Context, name string) ([]model.Manufacturer, error)
	CheckExists(ctx context.Context, manufacturerID int) (bool, error)
}

type ProductClassificationRepository interface {
	FindAll(ctx context.Context, name string) ([]model.ProductClassification, error)
	CheckExists(ctx context.Context, classificationID int) (bool, error)
}

type ProductFormRepository interface {
	FindAll(ctx context.Context, name string) ([]model.ProductForm, error)
	CheckExists(ctx context.Context, productFormID int) (bool, error)
}

type ProductRepository interface {
	CountAllByPharmacyID(ctx context.Context, pharmacyID int, filters map[string]interface{}) (int, error)
	FindAllByPharmacyID(ctx context.Context, pharmacyID int, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Product, error)
	InsertItem(ctx context.Context, product model.Product) error
	CheckExists(ctx context.Context, name string, genericName string, manufacturerId int) (bool, error)
	FindAll(ctx context.Context, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.Product, error)
	CountAll(ctx context.Context, filters map[string]interface{}) (int, error)
	IncrementUsage(ctx context.Context, id, diff int) error
	FindByID(ctx context.Context, id int) (*model.Product, error)
	FindIDBestseller(ctx context.Context, limit int) ([]int, error)
	SearchID(ctx context.Context, search string, offset, limit int) ([]int, error)
}

type PharmacyProductRepository interface {
	FindAllBestseller(ctx context.Context, bestSellerID []int) ([]model.Item, error)
	FindAllBestsellerPharmacies(ctx context.Context, productID, limit int) ([]model.AvailablePharmacy, error)
	FindAllRecommended(ctx context.Context, productIDs []int, minLat, maxLat, minLong, maxLong float64) ([]model.Item, error)
	Search(ctx context.Context, offset, limit, maxDistance int, search string, latitude, longitude decimal.Decimal) ([]model.Item, error)
	FindAllByPharmacyID(ctx context.Context, pharmacyID int, sortBy, sort string, offset, limit int, filters map[string]interface{}) ([]model.PharmacyProduct, error)
	Find(ctx context.Context, pharmacyID, productID int) (*model.PharmacyProduct, error)
	CountAllByPharmacyID(ctx context.Context, pharmacyID int, filters map[string]interface{}) (int, error)
	CountSearch(ctx context.Context, maxDistance int, search string, latitude, longitude decimal.Decimal) (int, error)
	Insert(ctx context.Context, pharmacyID int, item model.PharmacyProduct) error
	Update(ctx context.Context, pharmacyID, productID int, item model.PharmacyProduct) error
	UpdateStock(ctx context.Context, pharmacyID, productID, diff int) error
	Delete(ctx context.Context, pharmacyID, productID int) error
	CheckExists(ctx context.Context, pharmacyID, productID int) (bool, error)
	CheckUpdatedToday(ctx context.Context, pharmacyID, productID int) (bool, error)
	CheckSold(ctx context.Context, pharmacyID, productID int) (bool, error)
	FanOutFanInBestSeller(ctx context.Context, data []int, worker int) ([]model.Item, error)
	FanOutFanInRecommended(ctx context.Context, data []int, minLat, maxLat, minLong, maxLong float64, worker int) ([]model.Item, error)
}

type CartRepository interface {
	FindAll(ctx context.Context, userID int64) ([]model.CartItem, error)
	Insert(ctx context.Context, userID int64, pharmacyID, productID int, quantity int) error
	Update(ctx context.Context, userID int64, pharmacyID, productID, quantity int) error
	Delete(ctx context.Context, userID int64, pharmacyID, productID int) error
	DeleteOrdered(ctx context.Context, userID int64, pharmacies []int, products []int) error
	CheckExists(ctx context.Context, userID int64, pharmacyID, productID int) (bool, error)
}

type OrderRepository interface {
	Create(ctx context.Context, paymentID int, order model.Order) (int, error)
	FindAllUnpaid(ctx context.Context, userID int64) ([]model.UnpaidOrder, error)
	FindAll(ctx context.Context, userID int64, status string) ([]model.Order, error)
	FindAllByPharmacyID(ctx context.Context, pharmacyID int, status, sortBy, sort string, limit, offset int) ([]model.Order, error)
	CountAllByPharmacyID(ctx context.Context, pharmacyID int, status string) (int, error)
	Find(ctx context.Context, id int) (*model.Order, error)
	UpdateByOrderID(ctx context.Context, orderID int, status string) error
	UpdateByPaymentID(ctx context.Context, paymentID int, status string) error
	CheckExists(ctx context.Context, id int64, orderID int64) (bool, error)
}

type OrderItemRepository interface {
	BulkInsert(ctx context.Context, items []model.OrderItem) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment model.Payment) (int, error)
	Save(ctx context.Context, id int, imageURL string) error
}

type PharmacyLogisticRepository interface {
	FindAllByPharmacyID(ctx context.Context, addressID, pharmacyID int) ([]model.PharmacyLogistic, error)
	CheckEDA(ctx context.Context, logisticID, pharmacyID int) (int, error)
	Create(ctx context.Context, logisticID, pharmacyID int) error
}

type LogisticRepository interface {
	FindAll(ctx context.Context) ([]model.Logistic, error)
}
