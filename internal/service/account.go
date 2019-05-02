package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dchest/passwordreset"
	"github.com/google/uuid"
	"github.com/hadv/go-charity-me/internal/model"
	"github.com/hadv/go-charity-me/internal/repo"
	"github.com/pkg/errors"
)

// AccountService account service interface
type AccountService interface {
	Login(ctx context.Context, email, password string) (*model.User, error)
	Logout(ctx context.Context, email string) error
	Register(ctx context.Context, user *model.User) (*model.User, error)
	PasswordResetToken(ctx context.Context, email string) (string, error)
}

// Account service
type Account struct {
	repo repo.UserRepo
}

// NewAccount create new account service
func NewAccount(repo repo.UserRepo) *Account {
	return &Account{
		repo: repo,
	}
}

// Login check account login and generate user token
func (a *Account) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := a.repo.GetByEmail(ctx, email)
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
	usr, err := a.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (a *Account) Logout(ctx context.Context, email string) error {
	user, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	user.Token = ""
	_, err = a.repo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (a *Account) getPasswordHash(email string) ([]byte, error) {
	user, err := a.repo.GetByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return []byte(user.Password), nil
}

func (a *Account) PasswordResetToken(ctx context.Context, email string) (string, error) {
	user, err := a.repo.GetByEmail(context.Background(), email)
	if err != nil {
		return "", err
	}
	token := passwordreset.NewToken(email, 24*time.Hour, []byte(user.Password), signingKey)
	subject := "[Charity Me] Reset Password confirmation"
	plainTextContent := "Please click below link to reset your password: " + token
	htmlContent := fmt.Sprintf(`
		<strong>Please click below link to reset your password: </strong>
		<a href="http://localhost:8081/reset-password?token=%s">Reset Password</a>`, token)
	go sendMail(user, subject, plainTextContent, htmlContent)
	return token, nil
}

// Register create new account
func (a *Account) Register(ctx context.Context, user *model.User) (*model.User, error) {
	users, err := a.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, errors.New("email is already registered")
	}
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

	user, err = a.repo.Create(ctx, usr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot register new user")
	}

	return user, nil
}
