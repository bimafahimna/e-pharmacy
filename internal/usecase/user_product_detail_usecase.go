package usecase

import (
	"context"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

type UserProductDetailUseCase interface {
	GetProductDetails(ctx context.Context, uri dto.GetProductDetailURI, queries dto.GetProductDetailParams) (*dto.GetProductDetail, error)
	AvailablePharmacy(ctx context.Context, productID int) ([]dto.AvailablePharmacy, error)
}

func (u *userUseCase) GetProductDetails(ctx context.Context, uri dto.GetProductDetailURI, queries dto.GetProductDetailParams) (*dto.GetProductDetail, error) {
	product, err := u.store.Product().FindByID(ctx, uri.ProductID)
	if err != nil {
		return nil, err
	}

	pharmacy, err := u.store.Pharmacy().FindByID(ctx, int64(uri.PharmacyID))
	if err != nil {
		return nil, err
	}

	pharmacyProduct, err := u.store.PharmacyProduct().Find(ctx, uri.PharmacyID, uri.ProductID)
	if err != nil {
		return nil, err
	}

	pharmacist, err := u.store.Pharmacist().Find(ctx, int64(*pharmacy.PharmacistId))
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return nil, err
	}
	return &dto.GetProductDetail{
		ProductID:                     uri.ProductID,
		ProductName:                   product.Name,
		ProductGenericName:            product.GenericName,
		ProductImageUrl:               product.ImageURL,
		ProductForm:                   product.ProductForm,
		ProductManufacturer:           product.Manufacturer,
		ProductClassification:         product.ProductClassification,
		ProductCategory:               product.Categories,
		ProductDescription:            product.Description,
		ProductSellingUnit:            product.SellingUnit,
		ProductUnitInPack:             product.UnitInPack,
		ProductWeight:                 product.Weight,
		SelectedProductStock:          pharmacyProduct.Stock,
		SelectedProductPrice:          pharmacyProduct.Price.String(),
		SelectedProductSoldAmount:     pharmacyProduct.SoldAmount,
		SelectedPharmacyID:            pharmacyProduct.PharmacyID,
		SelectedPharmacyName:          pharmacyProduct.Name,
		SelectedPharmacistName:        pharmacist.Name,
		SelectedPharmacistPhoneNumber: pharmacist.WhatsappNumber,
		SelectedPharmacistSIPANumber:  pharmacist.SipaNumber,
		SelectedPharmacyAddress:       pharmacy.Address,
		SelectedPharmacyCityName:      pharmacy.CityName,
	}, nil
}

func (u *userUseCase) AvailablePharmacy(ctx context.Context, productID int) ([]dto.AvailablePharmacy, error) {
	availPharmacy, err := u.store.PharmacyProduct().FindAllBestsellerPharmacies(ctx, productID, 5)
	if err != nil {
		return nil, err
	}
	res := make([]dto.AvailablePharmacy, len(availPharmacy))
	for i, pharma := range availPharmacy {
		res[i] = dto.AvailablePharmacy{
			PharmacyID:   pharma.PharmacyID,
			PharmacyName: pharma.PharmacyName,
			PartnerName:  pharma.PartnerName,
			PartnerLogo:  pharma.PartnerLogo,
			Address:      pharma.Address,
			CityName:     pharma.CityName,
			ProductID:    pharma.ProductID,
			ProductPrice: pharma.Price,
			Stock:        pharma.Stock,
		}
	}
	return res, nil
}
