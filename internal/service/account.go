package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hadv/go-charity-me/internal/model"
	"github.com/hadv/go-charity-me/internal/repo"
	"github.com/pkg/errors"
)

// AccountService account service interface
type AccountService interface {
	Login(ctx context.Context, email, password string) (*model.User, error)
	Register(ctx context.Context, user *model.User) (*model.User, error)
}

type Account struct {
	repo repo.UserRepo
}

func NewAccount(repo repo.UserRepo) *Account {
	return &Account{
		repo: repo,
	}
}

func (a *Account) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !comparePasswords(user.Password, []byte(password)) {
		return nil, errors.New("passwords are not match")
	}
	token, err := generateToken(email, user.Password)
	if err != nil {
		return nil, err
	}
	user.Token = token
	if err = a.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a *Account) Register(ctx context.Context, user *model.User) (*model.User, error) {
	if user.Password != user.ConfirmPassword {
		return nil, errors.New("password and confirm password are not matched")
	}
	usr := &model.User{
		ID:        uuid.New().String(),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  hashAndSalt([]byte(user.Password)),
	}

	if err := a.repo.Create(ctx, usr); err != nil {
		return nil, errors.Wrap(err, "cannot register new user")
	}

	if usr, err := a.repo.Get(ctx, usr.ID); err != nil {
		return nil, errors.Wrapf(err, "cannot get user by id=%s", usr.ID)
	}
	return usr, nil
}
