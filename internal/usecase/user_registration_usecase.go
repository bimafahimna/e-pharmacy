package usecase

import (
	"context"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/model"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
	"github.com/hibiken/asynq"
	"google.golang.org/api/idtoken"
)

type UserRegistrationUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Verify(ctx context.Context, req dto.VerifyRequest) error
	GoogleRegister(ctx context.Context, req dto.GoogleRegisterRequest) error
	GoogleRegisterCallback(ctx context.Context, req dto.GoogleRegisterCallback) error
}

func (u *userUseCase) Register(ctx context.Context, req dto.RegisterRequest) error {
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
			if user.IsVerified {
				return apperror.ErrEmailHasBeenRegistered
			}
			return apperror.ErrVerifyRegisteredEmail
		}

		newUser := model.User{
			Role:         appconst.RoleCustomer,
			Email:        req.Email,
			PasswordHash: &passwordHash,
			IsVerified:   false,
		}
		userID, err := s.User().Create(ctx, newUser)
		if err != nil {
			return err
		}

		exists, err := s.VerificationToken().CheckExists(ctx, userID)
		if err != nil {
			return err
		}
		if exists {
			return apperror.ErrVerifyRegisteredEmail
		}

		token, err := util.GenerateRandomString(32)
		if err != nil {
			return apperror.ErrInternalServerError
		}

		if err := s.VerificationToken().Save(ctx, userID, token); err != nil {
			return err
		}

		opts := []asynq.Option{
			asynq.Queue(appconst.QueueCritical),
			asynq.MaxRetry(1),
		}
		return u.producer.ProduceVerificationTask(
			ctx,
			&mq.VerificationPayload{
				Email: req.Email,
				Token: token,
			},
			opts...,
		)
	})
}

func (u *userUseCase) Verify(ctx context.Context, req dto.VerifyRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		userID, err := s.VerificationToken().Find(ctx, req.Token)
		if err != nil {
			if errors.Is(err, apperror.ErrNotFound) {
				return apperror.ErrInvalidVerificationToken
			}
			return err
		}

		user, err := s.User().FindByID(ctx, userID)
		if err != nil {
			return err
		}

		err = u.bcrypt.CompareHashAndPassword(*user.PasswordHash, req.Password)
		if err != nil {
			return apperror.ErrPasswordDoNotMatch
		}

		if err := s.User().Verify(ctx, userID); err != nil {
			return err
		}

		return s.VerificationToken().Revoke(ctx, userID)
	})
}

func (u *userUseCase) GoogleRegister(ctx context.Context, req dto.GoogleRegisterRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		payload, err := idtoken.Validate(ctx, req.Credential, u.config.ClientID)
		if err != nil {
			return apperror.ErrInvalidToken
		}

		email, ok := payload.Claims["email"]
		if !ok {
			return apperror.ErrInvalidToken
		}

		exists, err := s.User().CheckExists(ctx, email.(string))
		if err != nil {
			return err
		}
		if exists {
			return apperror.ErrEmailHasBeenRegistered
		}

		user := model.User{
			Role:         appconst.RoleCustomer,
			Email:        email.(string),
			PasswordHash: nil,
			IsVerified:   true,
		}
		_, err = s.User().Create(ctx, user)
		return err
	})
}

func (u *userUseCase) GoogleRegisterCallback(ctx context.Context, req dto.GoogleRegisterCallback) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		exists, err := s.User().CheckExists(ctx, req.Email)
		if err != nil {
			return err
		}
		if exists {
			return apperror.ErrEmailHasBeenRegistered
		}

		user := model.User{
			Role:         appconst.RoleCustomer,
			Email:        req.Email,
			PasswordHash: nil,
			IsVerified:   true,
		}
		_, err = s.User().Create(ctx, user)
		return err
	})
}
