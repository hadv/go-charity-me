package service_test

import (
	"context"

	"github.com/hadv/go-charity-me/internal/model"
	"github.com/stretchr/testify/mock"
)

type UserMockedRepo struct {
	mock.Mock
}

func (u *UserMockedRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	args := u.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *UserMockedRepo) FindByEmail(ctx context.Context, email string) ([]model.User, error) {
	args := u.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (u *UserMockedRepo) Get(ctx context.Context, id string) (*model.User, error) {
	return nil, nil
}
func (u *UserMockedRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (u *UserMockedRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, nil
}
