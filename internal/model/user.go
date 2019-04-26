package model

import "time"

type User struct {
	ID              string    `json:"id"`
	Firstname       string    `json:"firstname"`
	Lastname        string    `json:"lastname"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	ConfirmPassword string    `db:"-" json:"-"`
	Token           string    `json:"token"`
	CreatedAt       time.Time `db:"created_at" json:"-"`
	UpdatedAt       time.Time `db:"updated_at" json:"-"`
}
