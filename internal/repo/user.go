package repo

import (
	"context"

	"github.com/hadv/go-charity-me/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	_, err := u.db.ExecContext(ctx, "INSERT INTO `users` (`id`, `firstname`, `lastname`, `email`, `password`, `token`) VALUES(?, ?, ?, ?, ?, ?)",
		user.ID, user.Firstname, user.Lastname, user.Email, user.Password, user.Token)
	return errors.Wrap(err, "cannot insert new user")
}

func (u *User) Get(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM `users` WHERE `id` = ?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find user with id=%s", id)
	}
	return &user, nil
}

func (u *User) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.db.GetContext(ctx, &user, "SELECT * FROM `users` WHERE `email` = ?", email)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find user with email=%s", email)
	}
	return &user, nil
}

// Update update user
func (u *User) Update(ctx context.Context, user *model.User) error {
	_, err := u.db.ExecContext(ctx, "UPDATE `users` SET `firstname`=?, `lastname`=?, `email`=?, `token`=? WHERE `id`=?",
		user.Firstname, user.Lastname, user.Email, user.Token, user.ID)
	return errors.Wrap(err, "cannot update user")
}
