package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
	"github.com/hibiken/asynq"
)

type PartnerManagementUseCase interface {
	AddPartner(ctx context.Context, req dto.AddPartnerRequest) error
	GetPartnerByID(ctx context.Context, id int) (*dto.PartnerResponse, error)
	EditPartner(ctx context.Context, req dto.EditPartnerRequest) error
	EditPartnerDaysAndHours(ctx context.Context, req dto.EditPartnerDaysAndHoursRequest) error
	ListPartners(ctx context.Context, params dto.ListPartnerParams) ([]dto.Partner, *dto.Pagination, error)
	RemovePartner()
}

func (u *adminUseCase) AddPartner(ctx context.Context, req dto.AddPartnerRequest) error {
	if err := util.ValidatePartnerRequest(req.ActiveDays, req.OperationalStart, req.OperationalStop); err != nil {
		return err
	}
	partner := dto.ConvertAddPartnerDtoToModel(req)

	return u.store.Atomic(ctx, func(s repository.Store) error {
		partnerRepository := s.Partner()

		if _, err := partnerRepository.FindByName(ctx, partner.Name); err == nil {
			return apperror.ErrPartnerNameHasBeenRegistered
		}
		if err := partnerRepository.Create(ctx, *partner); err != nil {
			return err
		}
		return nil
	})
}

func (u *adminUseCase) GetPartnerByID(ctx context.Context, id int) (*dto.PartnerResponse, error) {
	partner, err := u.store.Partner().FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.PartnerResponse{
		ID:               id,
		Name:             partner.Name,
		LogoUrl:          partner.LogoUrl,
		YearFounded:      partner.YearFounded,
		ActiveDays:       partner.ActiveDays,
		OperationalStart: partner.OperationalStart,
		OperationalStop:  partner.OperationalStop,
		IsActive:         strconv.FormatBool(partner.IsActive),
	}, nil
}

func (u *adminUseCase) EditPartner(ctx context.Context, req dto.EditPartnerRequest) error {
	if err := util.ValidatePartnerRequest(req.ActiveDays, req.OperationalStart, req.OperationalStop); err != nil {
		return err
	}
	partner := dto.ConvertEditPartnerDtoToModel(req)

	partnerRepository := u.store.Partner()
	if err := partnerRepository.UpdateProfile(ctx, *partner); err != nil {
		return err
	}

	opts := []asynq.Option{
		asynq.Queue(appconst.QueueCritical),
		asynq.MaxRetry(5),
		asynq.ProcessAt(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
	}
	return u.producer.ProduceUpdateDaysAndHours(ctx, &mq.UpdatePartnerDaysAndHoursPayload{
		ID:               req.ID,
		ActiveDays:       req.ActiveDays,
		OperationalStart: req.OperationalStart,
		OperationalStop:  req.OperationalStop,
	}, opts...)
}

func (u *adminUseCase) EditPartnerDaysAndHours(ctx context.Context, req dto.EditPartnerDaysAndHoursRequest) error {
	if err := util.ValidatePartnerRequest(req.ActiveDays, req.OperationalStart, req.OperationalStop); err != nil {
		return err
	}
	return u.store.Partner().UpdateDaysAndHours(ctx, model.Partner{
		ID:               int64(req.ID),
		ActiveDays:       util.IntArrayToString(req.ActiveDays),
		OperationalStart: req.OperationalStart,
		OperationalStop:  req.OperationalStop,
	})
}

func (u *adminUseCase) ListPartners(ctx context.Context, params dto.ListPartnerParams) ([]dto.Partner, *dto.Pagination, error) {
	repository := u.store.Partner()
	offset := (params.Page - 1) * params.Limit

	filters := map[string]interface{}{
		"name":        params.Name,
		"active_days": params.ActiveDays,
	}

	if params.Id != nil {
		filters["id"] = params.Id
	}

	if params.YearFounded != nil {
		filters["year_founded"] = params.YearFounded
	}

	if params.OperationalStart != nil {
		filters["operational_start"] = params.OperationalStart
	}

	if params.OperationalStop != nil {
		filters["operational_Stop"] = params.OperationalStop
	}

	if params.IsActive != nil {
		isAssigned, _ := strconv.ParseBool(*params.IsActive)
		filters["is_active"] = isAssigned
	}

	partners, err := repository.FindAll(ctx, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.Partner{}
	for _, p := range partners {
		resp = append(resp, dto.Partner{
			ID:               p.ID,
			Name:             p.Name,
			LogoUrl:          p.LogoUrl,
			YearFounded:      p.YearFounded,
			ActiveDays:       p.ActiveDays,
			OperationalStart: p.OperationalStart,
			OperationalStop:  p.OperationalStop,
			IsActive:         p.IsActive,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		})
	}

	totalRecords, err := repository.CountAll(ctx, filters)
	if err != nil {
		return nil, nil, err
	}
	return resp, dto.PaginationInfo(totalRecords, offset, params.Limit), nil
}

func (u *adminUseCase) RemovePartner() {
	panic("not implemented") // TODO: Implement
}
