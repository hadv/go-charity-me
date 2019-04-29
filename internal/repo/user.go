package repo

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/hadv/go-charity-me/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
}

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{
		db: db,
	}
}

type CircuitBreakerUser struct {
	User UserRepo
}

func (c *CircuitBreakerUser) Create(ctx context.Context, user *model.User) (*model.User, error) {
	output := make(chan *model.User, 1)
	hystrix.ConfigureCommand("create_user", hystrix.CommandConfig{
		Timeout: 1000,
	})
	errors := hystrix.Go("create_user", func() error {
		usr, err := c.User.Create(ctx, user)
		if err != nil {
			return err
		}
		output <- usr

		return nil
	}, nil)

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		return nil, err
	}
}

func (c *CircuitBreakerUser) Get(ctx context.Context, id string) (*model.User, error) {
	output := make(chan *model.User, 1)
	hystrix.ConfigureCommand("get_user_by_id", hystrix.CommandConfig{
		Timeout: 1000,
	})
	errors := hystrix.Go("get_user_by_id", func() error {
		user, err := c.User.Get(ctx, id)
		if err != nil {
			return err
		}
		output <- user

		return nil
	}, nil)

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		return nil, err
	}
}

func (c *CircuitBreakerUser) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	output := make(chan *model.User, 1)
	hystrix.ConfigureCommand("get_user_by_email", hystrix.CommandConfig{
		Timeout: 10000,
	})
	errors := hystrix.GoC(ctx, "get_user_by_email", func(ctx context.Context) error {
		user, err := c.User.GetByEmail(ctx, email)
		if err != nil {
			return err
		}
		output <- user

		return nil
	}, func(ctx context.Context, err error) error {
		fmt.Println(err)
		return err
	})

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		fmt.Println(err)
		return nil, err
	}
}

func (c *CircuitBreakerUser) Update(ctx context.Context, user *model.User) (*model.User, error) {
	output := make(chan *model.User, 1)
	hystrix.ConfigureCommand("update_user", hystrix.CommandConfig{
		Timeout: 1000,
	})
	errors := hystrix.Go("update_user", func() error {
		usr, err := c.User.Update(ctx, user)
		if err != nil {
			return err
		}
		output <- usr

		return nil
	}, nil)

	select {
	case out := <-output:
		return out, nil
	case err := <-errors:
		return nil, err
	}
}

func (u *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := u.db.ExecContext(ctx, "INSERT INTO `users` (`id`, `firstname`, `lastname`, `email`, `password`, `token`) VALUES(?, ?, ?, ?, ?, ?)",
		user.ID, user.Firstname, user.Lastname, user.Email, user.Password, user.Token)

	var usr model.User
	err = u.db.GetContext(ctx, &usr, "SELECT * FROM `users` WHERE `id` = ?", user.ID)
	if err != nil {
		return nil, err
	}
	return &usr, err
}

func (u *User) Get(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM `users` WHERE `id` = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM `users` WHERE `email` = ?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update update user
func (u *User) Update(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := u.db.ExecContext(ctx, "UPDATE `users` SET `firstname`=?, `lastname`=?, `email`=?, `token`=? WHERE `id`=?",
		user.Firstname, user.Lastname, user.Email, user.Token, user.ID)

	var usr model.User
	err = u.db.GetContext(ctx, &usr, "SELECT * FROM `users` WHERE `id` = ?", user.ID)
	if err != nil {
		return nil, err
	}
	return &usr, err
}
