package repo

import (
	"context"

	"github.com/hadv/go-charity-me/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
)

// UserRepo User repository interface
type UserRepo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) ([]model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
}

// User repository
type User struct {
	cb *gobreaker.CircuitBreaker
	db *sqlx.DB
}

// NewUser create new user repository
func NewUser(cb *gobreaker.CircuitBreaker, db *sqlx.DB) *User {
	return &User{
		cb: cb,
		db: db,
	}
}

// Create insert new user into db
func (u *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	usr, err := u.cb.Execute(func() (interface{}, error) {
		_, err := u.db.ExecContext(ctx, "INSERT INTO `users` (`id`, `firstname`, `lastname`, `email`, `password`, `token`) VALUES(?, ?, ?, ?, ?, ?)",
			user.ID, user.Firstname, user.Lastname, user.Email, user.Password, user.Token)
		if err != nil {
			return nil, errors.Wrap(err, "cannot insert new user")
		}
		var usr model.User
		err = u.db.GetContext(ctx, &usr, "SELECT * FROM `users` WHERE `id` = ?", user.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot find user with id=%s", user.ID)
		}

		return &usr, nil
	})

	if err != nil {
		return nil, err
	}
	return usr.(*model.User), nil
}

// Get Get user by id
func (u *User) Get(ctx context.Context, id string) (*model.User, error) {
	usr, err := u.cb.Execute(func() (interface{}, error) {
		var user model.User
		err := u.db.GetContext(ctx, &user, "SELECT * FROM `users` WHERE `id` = ?", id)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot find user with id=%s", id)
		}
		return &user, nil
	})

	if err != nil {
		return nil, err
	}
	return usr.(*model.User), nil
}

// GetByEmail Get user by email
func (u *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	usr, err := u.cb.Execute(func() (interface{}, error) {
		var user model.User
		err := u.db.GetContext(ctx, &user, "SELECT * FROM `users` WHERE `email` = ?", email)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot find user with email=%s", email)
		}
		return &user, nil
	})

	if err != nil {
		return nil, err
	}
	return usr.(*model.User), nil
}

// FindByEmail list users by email
func (u *User) FindByEmail(ctx context.Context, email string) ([]model.User, error) {
	list, err := u.cb.Execute(func() (interface{}, error) {
		var user []model.User
		err := u.db.SelectContext(ctx, &user, "SELECT * FROM `users` WHERE `email` = ?", email)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot find user with email=%s", email)
		}
		return user, nil
	})

	if err != nil {
		return nil, err
	}
	return list.([]model.User), nil
}

// Update update user
func (u *User) Update(ctx context.Context, user *model.User) (*model.User, error) {
	usr, err := u.cb.Execute(func() (interface{}, error) {
		_, err := u.db.ExecContext(ctx, "UPDATE `users` SET `firstname`=?, `lastname`=?, `email`=?, `password`=?, `token`=? WHERE `id`=?",
			user.Firstname, user.Lastname, user.Email, user.Password, user.Token, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "cannot update user")
		}
		var usr model.User
		err = u.db.GetContext(ctx, &usr, "SELECT * FROM `users` WHERE `id` = ?", user.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot find user with id=%s", user.ID)
		}
		return &usr, nil
	})

	if err != nil {
		return nil, err
	}
	return usr.(*model.User), nil
}
