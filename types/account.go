package types 

import (
	"github.com/google/uuid"
	"math/rand"
)

type Account struct{
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Number    int64   `json:"number"`
	Balance   int64   `json:"balance"`
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
	}
}