package handler

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hadv/go-charity-me/internal/model"
	"github.com/hadv/go-charity-me/internal/service"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Bind(r *http.Request) error {
	return nil
}

type RegisterRequest struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (l *RegisterRequest) Bind(r *http.Request) error {
	return nil
}

// Account handler all http request related to account
type Account struct {
	account service.AccountService
}

// NewAccount create new account handler
func NewAccount(account service.AccountService) *Account {
	return &Account{
		account: account,
	}
}

func (a *Account) Login(w http.ResponseWriter, r *http.Request) {
	login := &LoginRequest{}
	if err := render.Bind(r, login); err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.account.Login(r.Context(), login.Email, login.Password)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseData(w, r, user)
}

func (a *Account) Register(w http.ResponseWriter, r *http.Request) {
	registar := &RegisterRequest{}
	if err := render.Bind(r, registar); err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user := &model.User{
		Firstname:       registar.Firstname,
		Lastname:        registar.Lastname,
		Email:           registar.Email,
		Password:        registar.Password,
		ConfirmPassword: registar.ConfirmPassword,
	}

	user, err := a.account.Register(r.Context(), user)
	if err != nil {
		responseError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	responseData(w, r, user)
}
