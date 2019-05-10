package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/hadv/go-charity-me/internal/model"
	"github.com/hadv/go-charity-me/internal/service"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccount_Register(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()

	t.Run("error on finding user by email", func(t *testing.T) {
		user := &model.User{
			ID:    ID,
			Email: "dvietha@gmail.com",
		}
		repo := new(UserMockedRepo)
		repo.On("FindByEmail", ctx, user.Email).Return(nil, errors.New("cannot find user by email"))

		userService := service.NewAccount(repo)
		_, err := userService.Register(ctx, user)
		repo.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "cannot find user by email")
		}
	})

	t.Run("user already existed", func(t *testing.T) {
		user := &model.User{
			ID:    ID,
			Email: "dvietha@gmail.com",
		}
		repo := new(UserMockedRepo)
		repo.On("FindByEmail", ctx, user.Email).Return([]model.User{*user}, nil)
		userService := service.NewAccount(repo)
		_, err := userService.Register(ctx, user)
		repo.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "email is already registered")
		}
	})

	t.Run("passwords are not match", func(t *testing.T) {
		user := &model.User{
			ID:              ID,
			Email:           "dvietha@gmail.com",
			Password:        "12345678",
			ConfirmPassword: "87654321",
		}
		repo := new(UserMockedRepo)
		repo.On("FindByEmail", ctx, user.Email).Return([]model.User{}, nil)
		userService := service.NewAccount(repo)
		_, err := userService.Register(ctx, user)
		repo.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "password and confirm password are not matched")
		}
	})

	t.Run("error on creating new user", func(t *testing.T) {
		user := &model.User{
			ID:              ID,
			Email:           "dvietha@gmail.com",
			Password:        "12345678",
			ConfirmPassword: "12345678",
		}
		repo := new(UserMockedRepo)
		repo.On("FindByEmail", ctx, user.Email).Return([]model.User{}, nil)
		repo.On("Create", ctx, mock.Anything).Return(nil, errors.New("cannot create new user"))
		userService := service.NewAccount(repo)
		_, err := userService.Register(ctx, user)
		repo.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "cannot register new user: cannot create new user")
		}
	})

	t.Run("success on creating new user", func(t *testing.T) {
		user := &model.User{
			ID:              ID,
			Email:           "dvietha@gmail.com",
			Password:        "12345678",
			ConfirmPassword: "12345678",
		}
		repo := new(UserMockedRepo)
		repo.On("FindByEmail", ctx, user.Email).Return([]model.User{}, nil)
		repo.On("Create", ctx, mock.Anything).Return(user, nil)
		userService := service.NewAccount(repo)
		actual, err := userService.Register(ctx, user)
		repo.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, user, actual)
	})
}
