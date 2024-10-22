package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/util"
)

type UserManagementUseCase interface {
	ListUsers(ctx context.Context, params dto.ListUserParams) ([]dto.UserResponse, *dto.Pagination, error)
}

func (u *adminUseCase) ListUsers(ctx context.Context, params dto.ListUserParams) ([]dto.UserResponse, *dto.Pagination, error) {
	offset := util.ToOffset(params.Page, params.Limit)
	filters := params.Filters()

	userRepository := u.store.User()
	users, err := userRepository.FindAll(ctx, params.SortBy, params.Sort, offset, params.Limit, filters)
	if err != nil {
		return nil, nil, err
	}

	res := dto.ConvertToListUsersResponse(users)
	totalRecords, err := userRepository.CountAll(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	return res, dto.PaginationInfo(totalRecords, offset, params.Limit), nil
}
