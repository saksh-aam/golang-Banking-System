package main

import (
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ApiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct{
	Error string `json:"error"`
}

type CreateAccountRequest struct{
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}
type Account struct {
	ID        int `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64 `json:"number"`
	EncryptedPassword string `json:"-"`
	Balance   int64 `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}
type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	FirstName   string  `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}

type AddFundRequest struct {
	Amount    int `json:"amount"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}
type TransferReponse struct {
	Message string `json:"message"`
}

func (a *Account) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1000000)),
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}