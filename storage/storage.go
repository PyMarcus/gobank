package storage

import (
	"database/sql"
	"log"

	"github.com/PyMarcus/gobank/types"
	_ "github.com/lib/pq"
)

type Storage interface{
	CreateAccount(*types.Account) error 
	DeleteAccount(string) error 
	UpdateAccount(*types.Account) error 
	GetAccountById(string) (*types.Account, error)
	GetAccount() ([]*types.Account, error)
}

type PostgresqlStore struct{
	db *sql.DB 
}

func NewPostgresqlStore() (*PostgresqlStore, error){
	const CONNSTR string = "user=postgres dbname=fusion password=123 host=localhost port=5432 sslmode=disable"

	db, err := sql.Open("postgres", CONNSTR)

	if err != nil{
		log.Fatal(err)
	}

	log.Println("Connected into database.")

	return &PostgresqlStore{
		db: db,
	}, nil
}

func (psql *PostgresqlStore) CreateAccount(account *types.Account) error{
	log.Println("Inserting data...")
	stmt := "INSERT INTO account (id, first_name, last_name, number, balance, create_at) VALUES ($1, $2, $3, $4, $5, $6);"
	tx, err := psql.db.Begin()

	if err != nil {
		log.Panic(err)
	}

	// execute the insert
	_, err = tx.Exec(
		stmt,
		account.ID,
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

func (psql *PostgresqlStore) DeleteAccount(id string) error{
	return nil
}

func (psql *PostgresqlStore) UpdateAccount(account *types.Account) error{
	return nil
}

func (psql *PostgresqlStore) GetAccountById(id string) (*types.Account, error){
	return nil, nil
}

func (psql *PostgresqlStore) GetAccount() ([]*types.Account, error){
	accs := []*types.Account{}

	stmt := "SELECT * FROM account;"

	// execute the insert
	response, err := psql.db.Query(stmt)

	if err != nil {
		log.Panic(err)
	}

	defer response.Close()

	for response.Next(){
		acc := new(types.Account)

		response.Scan(
			&acc.ID,
			&acc.FirstName,
			&acc.LastName,
			&acc.Number,
			&acc.Balance,
			&acc.CreateAt,
		)

		accs = append(accs, acc)
	}

	log.Println(accs)
	return accs, nil
}
