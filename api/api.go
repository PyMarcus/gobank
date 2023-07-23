package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PyMarcus/gobank/storage"
	"github.com/PyMarcus/gobank/types"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

type APIError struct {
	Error string
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (a *APIServer) Run() {
	router := mux.NewRouter()
	
	router.HandleFunc("/account", makeHTTPHandleFunc(a.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(a.handleAccountId), a.store))
	router.HandleFunc("/transfer", withJWTAuth(makeHTTPHandleFunc(a.handleTransfer), a.store))
	log.Println("GOBANK API is running on ", a.listenAddr)
	http.ListenAndServe(a.listenAddr, router)
}

func (a *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	log.Println("Method: ", r.Method)
	switch r.Method {
	case "GET":
		return a.handleGetAccount(w, r)
	case "POST":
		return a.handleCreateAccount(w, r)
	case "DELETE":
		return a.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method %s not allowed!", r.Method)
}

func (a *APIServer) handleAccountId(w http.ResponseWriter, r *http.Request) error{
	log.Println("Method: ", r.Method)
	switch r.Method {
	case "GET":
		return a.handleGetAccount(w, r)
	case "DELETE":
		return a.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method %s not allowed!", r.Method)
}

func (a *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	resp, err := a.store.GetAccount()
	if err != nil {
		log.Panicln(err)
	}
	WriteJSON(w, http.StatusOK, resp)
	log.Println("Successfully get data")
	return nil
}

func (a *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	account, err := a.store.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (a *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	create := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(create); err != nil {
		return err
	}

	acc := types.NewAccount(create.FirstName, create.LastName)

	if err := a.store.CreateAccount(acc); err != nil {
		return err
	}

	token, err := createJWTToken(acc)

	if err != nil{
		return err 
	}

	log.Println(token)

	return WriteJSON(w, http.StatusCreated, acc)
}

func (a *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	log.Println("Deleted method")
	account, err := a.store.GetAccountById(id)
	if err != nil {
		return err
	}

	err = a.store.DeleteAccount(id)

	if err == nil {
		return WriteJSON(w, http.StatusOK, account)
	}
	return err
}

func (a *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	toTransfer := new(types.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(toTransfer); err != nil{
		return err 
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, toTransfer)
}

func createJWTToken(account *types.Account) (string, error){
	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(time.Hour * 24).Unix(),
		"accountNumber": account.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func withJWTAuth(handleFunc http.HandlerFunc, s storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		log.Println("Calling JWT middleware to protect route...")

		tokenStr := r.Header.Get("X-token")
		
		tk, err := validateJWT(tokenStr)

		if !tk.Valid{
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return 
		}

		if err != nil{
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return 
		}


		if err != nil{
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
		}else{
			handleFunc(w, r)
		}
	}
}

func getID(r *http.Request) string{
	return mux.Vars(r)["id"]
}

const tok = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo1MjE1MTQsImV4cGlyZXNBdCI6MTY5MDIyOTMwMn0.ynFoS5SgzSv-8JY3R6_xRSQq96E3YwGPk9EOE-sjnuY"

func validateJWT(tokenStr string) (*jwt.Token, error){
	JWT_SECRET := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(JWT_SECRET), nil
	})
	
}

type apiFunc func(http.ResponseWriter, *http.Request) error

// makeHTTPHandleFunc is a decorator to a handle function to convert type
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, values any) error {
	w.Header().Set("Content-Type", "application/json") //this must be come before the writterheader function
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(values)
}
