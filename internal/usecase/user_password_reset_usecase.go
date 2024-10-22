package usecase

import (
	"context"
	"errors"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/appconst"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
	"github.com/hibiken/asynq"
)

type UserPasswordResetUseCase interface {
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	ConfirmResetPassword(ctx context.Context, req dto.ConfirmResetRequest) error
}

func (u *userUseCase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	return u.store.Atomic(ctx, func(s repository.Store) error {
		user, err := s.User().FindByEmail(ctx, req.Email)
		if err != nil {
			return err
		}

		if user.PasswordHash == nil {
			return apperror.ErrCannotResetPassword
		}

		tokenRepository := s.PasswordResetToken()
		ok, err := tokenRepository.CheckExists(ctx, user.ID)
		if err != nil {
			return err
		}
		if ok {
			return apperror.ErrResetPasswordTokenExists
		}

		token, err := util.GenerateRandomString(32)
		if err != nil {
			return apperror.ErrInternalServerError
		}

		if err := tokenRepository.Save(ctx, user.ID, token); err != nil {
			return err
		}

		opts := []asynq.Option{
			asynq.Queue(appconst.QueueCritical),
			asynq.MaxRetry(1),
		}
		return u.producer.ProducePasswordResetTask(
			ctx,
			&mq.PasswordResetPayload{
				Email: req.Email,
				Token: token,
			},
			opts...,
		)
	})
}

func (u *userUseCase) ConfirmResetPassword(ctx context.Context, req dto.ConfirmResetRequest) error {
	if err := util.ValidatePassword(req.Password); err != nil {
		return apperror.ErrInvalidPasswordFormat
	}

	passwordHash, err := u.bcrypt.Hash(req.Password)
	if err != nil {
		return apperror.ErrInternalServerError
	}

	return u.store.Atomic(ctx, func(s repository.Store) error {
		userRepository := s.User()
		tokenRepository := s.PasswordResetToken()

		userID, err := tokenRepository.FindUserID(ctx, req.Token)
		if err != nil {
			if errors.Is(err, apperror.ErrNotFound) {
				return apperror.ErrInvalidPasswordResetToken
			}
			return err
		}

		if err = userRepository.UpdatePassword(ctx, userID, passwordHash); err != nil {
			return err
		}

		if err := tokenRepository.Revoke(ctx, userID); err != nil {
			return err
		}

		return nil
	})
}
