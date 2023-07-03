package types 

import (
	"time"
	"github.com/google/uuid"
	"math/rand"
)

type Account struct{
	ID        string    `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Number    int64     `json:"number" db:"number"`
	Balance   int64     `json:"balance" db:"balance"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
}

type CreateAccountRequest struct{
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
}

func createUUID() string{
	id := uuid.New()
	return id.String()
}

func NewAccount(firstName, lastName string) *Account{
	return &Account{
		ID: createUUID(),
		FirstName: firstName,
		LastName: lastName,
		Number: rand.Int63n(1000000),
		Balance: rand.Int63n(100000),
		CreateAt: time.Now().UTC(),
	}
}