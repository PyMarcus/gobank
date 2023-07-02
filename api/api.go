package api

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/gorilla/mux"
)

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
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(values)
}

type APIServer struct {
	listenAddr string
}

type APIError struct{
	Error string 
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (a *APIServer) Run(){
	router := mux.NewRouter()
	
	router.HandleFunc("/account", makeHTTPHandleFunc(a.handleAccount))

	log.Println("GOBANK API is running on ", a.listenAddr)
	http.ListenAndServe(a.listenAddr, router)
}

func (a *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error{
	return nil 
}

func (a *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error{
	return nil 
}

func (a *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error{
	return nil 
}

func (a *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error{
	return nil 
}

func (a *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error{
	return nil 
}
