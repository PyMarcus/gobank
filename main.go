package main

import (
	"log"

	"github.com/PyMarcus/gobank/api"
	"github.com/PyMarcus/gobank/storage"
)

func main(){

	store, err := storage.NewPostgresqlStore()

	if err != nil{
		log.Fatal(err)
	}
	log.Println("running...")
	a := api.NewAPIServer("127.0.0.1:3000", store)
	a.Run()
	log.Println("Ok!")
}
