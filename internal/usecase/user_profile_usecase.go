package usecase

import (
	"context"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
)

type UserProfileUseCase interface {
	UpdateProfile()
	AddAddress(ctx context.Context, req dto.CustomerAddressRequest, userId int64) error
	GetAddresses(ctx context.Context, userId int64) ([]dto.CustomerAddressResponse, error)
}

func (u *userUseCase) UpdateProfile() {
	panic("not implemented") // TODO: Implement
}

func (u *userUseCase) AddAddress(ctx context.Context, req dto.CustomerAddressRequest, userId int64) error {
	newAddress := dto.ConvertToCustomerAddressModel(req)
	newAddress.UserID = userId
	return u.store.Atomic(ctx, func(s repository.Store) error {
		customerAddressRepository := s.CustomerAddress()
		if newAddress.IsActive {
			address, err := customerAddressRepository.FindActiveByUserId(ctx, userId)
			if err != nil {
				return err
			}
			if address != nil {
				if err := customerAddressRepository.UnsetActive(ctx, address.ID); err != nil {
					return err
				}
			}
		}
		return customerAddressRepository.Create(ctx, *newAddress)
	})
}

func (u *userUseCase) GetAddresses(ctx context.Context, userId int64) ([]dto.CustomerAddressResponse, error) {
	res := []dto.CustomerAddressResponse{}
	customerAddressRepository := u.store.CustomerAddress()
	addresses, err := customerAddressRepository.FindAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		address.UserID = userId
		res = append(res, *dto.ConvertToCustomerAddressResponseDto(address))
	}
	return res, nil
}
