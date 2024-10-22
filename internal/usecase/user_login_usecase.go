package usecase

import (
	"context"
	"strconv"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"google.golang.org/api/idtoken"
)

type UserLoginUseCase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	GoogleLogin(ctx context.Context, req dto.GoogleLoginRequest) (*dto.LoginResponse, error)
	GoogleLoginCallback(ctx context.Context, req dto.GoogleLoginCallback) (*dto.LoginResponse, error)
}

func (u *userUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.store.User().FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// TODO: Ensure only users with role "customer" may perform login here

	if err = u.bcrypt.CompareHashAndPassword(*user.PasswordHash, req.Password); err != nil {
		return nil, apperror.ErrIncorrectEmailOrPassword
	}

	jwt, err := u.jwt.Sign(user.ID, user.Role, user.IsVerified)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}

	isVerified := strconv.FormatBool(user.IsVerified)
	return &dto.LoginResponse{
		Token:      jwt,
		Role:       user.Role,
		IsVerified: isVerified,
		UserID:     user.ID,
	}, nil
}

func (u *userUseCase) GoogleLogin(ctx context.Context, req dto.GoogleLoginRequest) (*dto.LoginResponse, error) {
	payload, err := idtoken.Validate(ctx, req.Credential, u.config.ClientID)
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	email, ok := payload.Claims["email"]
	if !ok {
		return nil, apperror.ErrInvalidToken
	}

	user, err := u.store.User().FindByEmail(ctx, email.(string))
	if err != nil {
		return nil, err
	}

	token, err := u.jwt.Sign(user.ID, user.Role, user.IsVerified)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}

	isVerified := strconv.FormatBool(user.IsVerified)
	return &dto.LoginResponse{
		Token:      token,
		Role:       user.Role,
		IsVerified: isVerified,
	}, nil
}

func (u *userUseCase) GoogleLoginCallback(ctx context.Context, req dto.GoogleLoginCallback) (*dto.LoginResponse, error) {
	user, err := u.store.User().FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	token, err := u.jwt.Sign(user.ID, user.Role, user.IsVerified)
	if err != nil {
		return nil, apperror.ErrInternalServerError
	}

	isVerified := strconv.FormatBool(user.IsVerified)
	return &dto.LoginResponse{
		UserID:     user.ID,
		Token:      token,
		Role:       user.Role,
		IsVerified: isVerified,
	}, nil
}
