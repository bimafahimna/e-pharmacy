package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
	"github.com/hibiken/asynq"
)

type PharmacistManagementUseCase interface {
	AddPharmacist(ctx context.Context, req dto.AddPharmacistRequest) error
	ListPharmacists(ctx context.Context, params dto.ListPharmacistParams) ([]dto.Pharmacist, *dto.Pagination, error)
	GetPharmacist(ctx context.Context, userID int64) (*dto.Pharmacist, error)
	EditPharmacist(ctx context.Context, userID int64, req dto.UpdatePharmacistRequest) error
	RemovePharmacist(ctx context.Context, userID int64) error
}

func (u *adminUseCase) AddPharmacist(ctx context.Context, req dto.AddPharmacistRequest) error {
	if err := util.ValidatePassword(req.Password); err != nil {
		return apperror.ErrInvalidPasswordFormat
	}

	passwordHash, err := u.bcrypt.Hash(req.Password)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	return u.store.Atomic(ctx, func(s repository.Store) error {
		user, err := s.User().FindByEmail(ctx, req.Email)
		if errors.Is(err, apperror.ErrInternalServerError) {
			return err
		}
		if user != nil {
			if user.DeletedAt == nil {
				return apperror.ErrEmailHasBeenRegistered
			}

			if err := s.User().Recover(ctx, user.ID); err != nil {
				return err
			}

			pharmacist := model.Pharmacist{
				UserID:            user.ID,
				Name:              req.Name,
				SipaNumber:        req.SipaNumber,
				WhatsappNumber:    req.WhatsappNumber,
				YearsOfExperience: req.YearsOfExperience,
				IsAssigned:        false,
			}
			err := s.Pharmacist().Create(ctx, pharmacist)
			if errors.Is(err, apperror.ErrNotFound) {
				return apperror.ErrPharmacistMustBeUnique
			}
			return err
		}

		newUser := model.User{
			Role:         appconst.RolePharmacist,
			Email:        req.Email,
			PasswordHash: &passwordHash,
			IsVerified:   true,
		}
		userID, err := s.User().Create(ctx, newUser)
		if err != nil {
			return err
		}

		pharmacist := model.Pharmacist{
			UserID:            userID,
			Name:              req.Name,
			SipaNumber:        req.SipaNumber,
			WhatsappNumber:    req.WhatsappNumber,
			YearsOfExperience: req.YearsOfExperience,
			IsAssigned:        false,
		}
		if err := s.Pharmacist().Create(ctx, pharmacist); err != nil {
			if errors.Is(err, apperror.ErrNotFound) {
				return apperror.ErrPharmacistMustBeUnique
			}
			return err
		}

		opts := []asynq.Option{
			asynq.Queue(appconst.QueueDefault),
			asynq.MaxRetry(1),
		}
		return u.producer.ProducePharmacistCredentialsTask(
			ctx,
			&mq.PharmacistCredentialsPayload{
				Email:    req.Email,
				Password: req.Password,
			},
			opts...,
		)
	})
}

func (u *adminUseCase) ListPharmacists(ctx context.Context, params dto.ListPharmacistParams) ([]dto.Pharmacist, *dto.Pagination, error) {
	offset := (params.Page - 1) * params.Limit

	filters := map[string]interface{}{
		"name":            params.Name,
		"email":           params.Email,
		"sipa_number":     params.SipaNumber,
		"whatsapp_number": params.WhatsappNumber,
	}

	if params.YearsOfExperience != "" {
		filters["years_of_experience"] = params.YearsOfExperience
	}

	if params.IsAssigned != "" {
		isAssigned, _ := strconv.ParseBool(params.IsAssigned)
		filters["is_assigned"] = isAssigned
	}

	pharmacists, err := u.store.Pharmacist().FindAll(ctx, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	resp := []dto.Pharmacist{}
	for _, p := range pharmacists {
		resp = append(resp, dto.Pharmacist{
			UserID:            p.UserID,
			Name:              p.Name,
			Email:             p.Email,
			SipaNumber:        p.SipaNumber,
			WhatsappNumber:    p.WhatsappNumber,
			YearsOfExperience: p.YearsOfExperience,
			IsAssigned:        p.IsAssigned,
			CreatedAt:         p.CreatedAt,
		})
	}

	totalRecords, err := u.store.Pharmacist().CountAll(ctx, filters)
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
	return resp, pagination, nil
}

func (u *adminUseCase) GetPharmacist(ctx context.Context, userID int64) (*dto.Pharmacist, error) {
	pharmacist, err := u.store.Pharmacist().Find(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &dto.Pharmacist{
		UserID:            pharmacist.UserID,
		Name:              pharmacist.Name,
		Email:             pharmacist.Email,
		SipaNumber:        pharmacist.SipaNumber,
		WhatsappNumber:    pharmacist.WhatsappNumber,
		YearsOfExperience: pharmacist.YearsOfExperience,
		IsAssigned:        pharmacist.IsAssigned,
	}
	return resp, nil
}

func (u *adminUseCase) EditPharmacist(ctx context.Context, userID int64, req dto.UpdatePharmacistRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		pharmacist, err := s.Pharmacist().Find(ctx, userID)
		if err != nil {
			return err
		}

		updatedPharmacist := model.Pharmacist{
			Name:              pharmacist.Name,
			SipaNumber:        pharmacist.SipaNumber,
			WhatsappNumber:    req.WhatsappNumber,
			YearsOfExperience: req.YearsOfExperience,
			IsAssigned:        pharmacist.IsAssigned,
		}
		if err := s.Pharmacist().Update(ctx, userID, updatedPharmacist); err != nil {
			if errors.Is(err, apperror.ErrUniqueViolation) {
				return apperror.ErrPharmacistMustBeUnique
			}
			return err
		}
		return nil
	})
}

func (u *adminUseCase) RemovePharmacist(ctx context.Context, userID int64) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		assigned, err := s.Pharmacist().CheckAssigned(ctx, userID)
		if err != nil {
			return err
		}
		if assigned {
			return apperror.ErrBadRequest
		}

		if err = s.User().Delete(ctx, userID); err != nil {
			return err
		}
		return s.Pharmacist().Delete(ctx, userID)
	})
}
