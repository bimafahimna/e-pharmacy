package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
)

type PharmacyManagementUseCase interface {
	AddPharmacy(ctx context.Context, req dto.AddPharmacyRequest) error
	ListPharmacies(ctx context.Context, param dto.ListPharmacyParams) ([]dto.PharmacyResponse, *dto.Pagination, error)
	UpdatePharmacy()
	RemovePharmacy()
	DownloadPharmacyMedicineStock()
}

func (u *adminUseCase) AddPharmacy(ctx context.Context, req dto.AddPharmacyRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		pharmacy := dto.ConvertToPharmacyModel(req)
		pharmacyID, err := s.Pharmacy().Create(ctx, *pharmacy)
		if err != nil {
			return err
		}
		for i := 0; i < len(req.Logistics); i++ {
			err := s.PharmacyLogistic().Create(ctx, req.Logistics[i], pharmacyID)
			if err != nil {
				return err
			}
		}
		pharmacist, err := s.Pharmacist().Find(ctx, int64(req.PharmacistID))
		if err != nil {
			return err
		}

		if pharmacist.IsAssigned {
			return apperror.ErrPharmacistIsAlreadyAssigned
		}
		pharmacist.IsAssigned = true
		return s.Pharmacist().Update(ctx, int64(req.PharmacistID), *pharmacist)
	})
}

func (u *adminUseCase) ListPharmacies(ctx context.Context, params dto.ListPharmacyParams) ([]dto.PharmacyResponse, *dto.Pagination, error) {
	repository := u.store.Pharmacy()
	offset := util.ToOffset(params.Page, params.Limit)

	filters := params.Filters()

	pharmacies, err := repository.FindAll(ctx, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	res := dto.ConvertToListPharmacies(pharmacies)
	totalRecords, err := repository.CountAll(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (totalRecords + params.Limit - 1) / params.Limit
	currentPage := offset/params.Limit + 1

	pagination := &dto.Pagination{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: currentPage > 1,
		NextPage:     currentPage < totalPages,
	}
	return res, pagination, nil
}

func (u *adminUseCase) UpdatePharmacy() {
	panic("not implemented") // TODO: Implement
}

func (u *adminUseCase) RemovePharmacy() {
	panic("not implemented") // TODO: Implement
}

func (u *adminUseCase) DownloadPharmacyMedicineStock() {
	panic("not implemented") // TODO: Implement
}
