package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/PyMarcus/gobank/types"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(string) error
	UpdateAccount(*types.Account) error
	GetAccountById(string) (*types.Account, error)
	GetAccount() ([]*types.Account, error)
	GetAccountByNumber(int64) (*types.Account, error)
}

type PostgresqlStore struct {
	db *sql.DB
}

func NewPostgresqlStore() (*PostgresqlStore, error) {
	const CONNSTR string = "user=postgres dbname=postgres host=localhost port=5432 sslmode=disable"

	db, err := sql.Open("postgres", CONNSTR)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected into database.")

	return &PostgresqlStore{
		db: db,
	}, nil
}

func (psql *PostgresqlStore) CreateAccount(account *types.Account) error {
	log.Println("Inserting data to ", *&account.FirstName)
	stmt := "INSERT INTO account (id, encrypted_password, first_name, last_name, number, balance, create_at) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	tx, err := psql.db.Begin()

	if err != nil {
		log.Panic(err)
	}

	// execute the insert
	_, err = tx.Exec(
		stmt,
		account.ID,
		account.EncryptedPassword,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreateAt,
	)

	if err != nil {
		tx.Rollback()
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	return nil
}

func (psql *PostgresqlStore) DeleteAccount(id string) error {
	_, err := psql.db.Exec("DELETE FROM account WHERE id = $1;", id)

	if err != nil {
		return err
	}

	log.Printf("Data with id: %s was deleted \n", id)
	return nil
}

func (psql *PostgresqlStore) GetAccountByNumber(number int64) (*types.Account, error){
	rows, err := psql.db.Query("SELECT * FROM account WHERE number = $1;", number)

	if err != nil {
		return nil, err
	}

	for rows.Next(){
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Not authorizated")
}

func (psql *PostgresqlStore) UpdateAccount(account *types.Account) error {
	acc, err := types.UpdateAccount(account.ID, account.FirstName, account.LastName, account.EncryptedPassword, account.Number, account.Balance)
	
	_, err = psql.db.Query("UPDATE account SET first_name = $1, last_name = $2, encrypted_password = $3, number = $4, balance = $5 WHERE id = $6;",
	acc.FirstName, acc.LastName, acc.EncryptedPassword, acc.Number, acc.Balance, acc.ID)

	if err != nil {
		return  err
	}

	return nil
}

func (psql *PostgresqlStore) GetAccountById(id string) (*types.Account, error) {
	rows, err := psql.db.Query("SELECT * FROM account WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account %s not found", id)
}

func (psql *PostgresqlStore) GetAccount() ([]*types.Account, error) {
	accs := []*types.Account{}

	stmt := "SELECT * FROM account;"

	// execute the insert
	response, err := psql.db.Query(stmt)

	if err != nil {
		log.Panic(err)
	}

	defer response.Close()

	for response.Next() {

		acc, err := scanIntoAccount(response)

		if err != nil {
			return nil, err
		}

		accs = append(accs, acc)
	}

	log.Println(accs)
	return accs, nil
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	acc := new(types.Account)

	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Balance,
		&acc.EncryptedPassword,
		&acc.CreateAt,
	)
	return acc, err
}
