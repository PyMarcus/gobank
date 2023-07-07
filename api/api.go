package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PyMarcus/gobank/storage"
	"github.com/PyMarcus/gobank/types"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store storage.Storage
}

type APIError struct{
	Error string 
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (a *APIServer) Run(){
	router := mux.NewRouter()
	
	router.HandleFunc("/account", makeHTTPHandleFunc(a.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(a.handleGetAccountById))
	
	log.Println("GOBANK API is running on ", a.listenAddr)
	http.ListenAndServe(a.listenAddr, router)
}

func (a *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error{
	switch(r.Method){
	case "GET":
		return a.handleGetAccount(w, r)
	case "POST":
		return a.handleCreateAccount(w, r)
	case "DELETE":
		return a.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method %s not allowed!", r.Method)
}

func (a *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error{
	resp, err := a.store.GetAccount()
	if err != nil{
		log.Panicln(err)
	}
	WriteJSON(w, http.StatusOK, resp)
	log.Println("Successfully get data")
	return nil
}

func (a *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error{
	id := mux.Vars(r)["id"]

	account, err := a.store.GetAccountById(id)

	if err != nil{
		return err 
	}

	return WriteJSON(w, http.StatusOK, account)
}


func (a *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error{
	create := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(create); err != nil{
		return err 
	}

	acc := types.NewAccount(create.FirstName, create.LastName)

	if err := a.store.CreateAccount(acc); err != nil{
		return err
	}

	return WriteJSON(w, http.StatusCreated, acc)
}

func (a *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error{
	return nil 
}

func (a *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error{
	return nil 
}


type apiFunc func(http.ResponseWriter, *http.Request) error

// makeHTTPHandleFunc is a decorator to a handle function to convert type
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w, r); err != nil{
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, values any) error{
	w.Header().Set("Content-Type", "application/json")   //this must be come before the writterheader function
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(values)
}
