package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type LocationUseCase interface {
	ListProvinces(ctx context.Context, name string) ([]dto.ProvinceResponse, error)
	ListCities(ctx context.Context, params dto.ListCityParams) ([]dto.CityResponse, error)
	ListDistricts(ctx context.Context, params dto.ListDistrictParams) ([]dto.DistrictResponse, error)
	ListSubDistricts(ctx context.Context, params dto.ListSubDistrictParams) ([]dto.SubDistrictResponse, error)
}

type locationUseCase struct {
	store repository.Store
}

func NewLocationUseCase(store repository.Store) LocationUseCase {
	return &locationUseCase{
		store: store,
	}
}

func (u *locationUseCase) ListProvinces(ctx context.Context, name string) ([]dto.ProvinceResponse, error) {
	provinceRepository := u.store.Province()
	provinces, err := provinceRepository.FindAll(ctx, name)
	if err != nil {
		return nil, err
	}

	res := []dto.ProvinceResponse{}
	for _, province := range provinces {
		res = append(res, dto.ProvinceResponse{
			ID:   province.ID,
			Name: province.Name,
		})
	}
	return res, nil
}

func (u *locationUseCase) ListCities(ctx context.Context, params dto.ListCityParams) ([]dto.CityResponse, error) {
	filters := map[string]interface{}{
		"name": params.Name,
	}

	if params.ProvinceId != nil {
		filters["province_unofficial_id"] = params.ProvinceId
	}

	cityRepository := u.store.City()
	cities, err := cityRepository.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	res := []dto.CityResponse{}
	for _, city := range cities {
		res = append(res, dto.CityResponse{
			ID:                   city.ID,
			ProvinceId:           city.ProvinceId,
			CityName:             city.Name,
			CityType:             city.Type,
			CityUnofficialId:     city.UnofficialId,
			ProvinceUnofficialId: city.ProvinceUnofficialId,
		})
	}
	return res, nil
}

func (u *locationUseCase) ListDistricts(ctx context.Context, params dto.ListDistrictParams) ([]dto.DistrictResponse, error) {
	filters := map[string]interface{}{
		"name": params.Name,
	}

	if params.CityId != nil {
		filters["city_id"] = params.CityId
	}

	districtRepository := u.store.District()
	districts, err := districtRepository.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	res := []dto.DistrictResponse{}
	for _, district := range districts {
		res = append(res, dto.DistrictResponse{
			ID:     district.ID,
			CityId: district.CityId,
			Name:   district.Name,
		})
	}
	return res, nil
}

func (u *locationUseCase) ListSubDistricts(ctx context.Context, params dto.ListSubDistrictParams) ([]dto.SubDistrictResponse, error) {
	filters := map[string]interface{}{
		"name": params.Name,
	}

	if params.DistrictId != nil {
		filters["district_id"] = params.DistrictId
	}

	subDistrictRepository := u.store.SubDistrict()
	districts, err := subDistrictRepository.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	res := []dto.SubDistrictResponse{}
	for _, district := range districts {
		res = append(res, dto.SubDistrictResponse{
			ID:         district.ID,
			DistrictId: district.DistrictId,
			Name:       district.Name,
		})
	}
	return res, nil
}
