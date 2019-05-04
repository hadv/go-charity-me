package model

import "time"

type User struct {
	ID                string    `db:"id" json:"id"`
	Firstname         string    `db:"firstname" json:"firstname"`
	Lastname          string    `db:"lastname" json:"lastname"`
	Email             string    `db:"email" json:"email"`
	EmailVerification bool      `db:"email_verification" json:"email_verification"`
	Password          string    `db:"password" json:"password"`
	ConfirmPassword   string    `db:"-" json:"-"`
	Token             string    `db:"token" json:"token"`
	CreatedAt         time.Time `db:"created_at" json:"-"`
	UpdatedAt         time.Time `db:"updated_at" json:"-"`
}
