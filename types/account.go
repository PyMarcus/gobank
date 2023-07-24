package types 

import (
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

type LoginRequest struct{
	Number int64    `json:"number" db:"number"`
	Password string `json:"password" db:"password"`
}

type LoginResponse struct{
	Token string    `json:"token"`
	Number int64	`json:"number"`
}

type Account struct{
	ID        		  string    `json:"id" db:"id"`
	EncryptedPassword string    `json:"-" db:"encrypted_password"`
	FirstName         string    `json:"first_name" db:"first_name"`
	LastName          string    `json:"last_name" db:"last_name"`
	Number            int64     `json:"number" db:"number"`
	Balance           int64     `json:"balance" db:"balance"`
	CreateAt          time.Time `json:"create_at" db:"create_at"`
}

type TransferRequest struct{
	ToAccount string `json:"to_account"`
	Amount int64 `json:"amount"`
}

type CreateAccountRequest struct{
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name"  db:"last_name"`
	Password  string    `json:"password"   db:"encrypted_password"`
}

func createUUID() string{
	id := uuid.New()
	return id.String()
}

func NewAccount(firstName, lastName, password string) (*Account, error){
	encryptedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)


	if err != nil{
		return nil, err 
	}

	return &Account{
		ID: createUUID(),
		FirstName: firstName,
		LastName: lastName,
		EncryptedPassword: string(encryptedPW),
		Number: rand.Int63n(1000000),
		Balance: rand.Int63n(100000),
		CreateAt: time.Now().UTC(),
	}, nil
}

func UpdateAccount(id, firstName, lastName, password string, number, balance int64) (*Account, error){
	encryptedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)


	if err != nil{
		return nil, err 
	}

	return &Account{
		ID: id,
		FirstName: firstName,
		LastName: lastName,
		EncryptedPassword: string(encryptedPW),
		Number: number,
		Balance: balance,
		CreateAt: time.Now().UTC(),
	}, nil
}
