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
	
	router.HandleFunc("/login", makeHTTPHandleFunc(a.handleLogin))
	router.HandleFunc("/account", withJWTAuth(makeHTTPHandleFunc(a.handleAccount), a.store))
	router.HandleFunc("/account/update", withJWTAuth(makeHTTPHandleFunc(a.handleUpdate), a.store))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(a.handleAccountId), a.store))
	router.HandleFunc("/transfer", withJWTAuth(makeHTTPHandleFunc(a.handleTransfer), a.store))
	log.Println("GOBANK API is running on ", a.listenAddr)
	http.ListenAndServe(a.listenAddr, router)
}

func (a *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error{
	if r.Method != "POST"{
		return fmt.Errorf("Method not allowed %s", r.Method)
	}

	var requestLogin types.LoginRequest

	
	if err := json.NewDecoder(r.Body).Decode(&requestLogin); err != nil{
		return err
	}

	acc, err := a.store.GetAccountByNumber(requestLogin.Number)

	if err != nil{
		return err 
	}

	if acc != nil{
		token, err := createJWTToken(acc)
		if err != nil{
			return err 
		}

		return WriteJSON(w, http.StatusOK, &types.LoginResponse{
			Token: token,
			Number: acc.Number,
		})
	}

	return WriteJSON(w, http.StatusOK, requestLogin)
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

func (a *APIServer) handleUpdate(w http.ResponseWriter, r *http.Request) error{
	if r.Method != "POST"{
		return fmt.Errorf("Method not allowed %s", r.Method)
	}

	var acc *types.Account

	json.NewDecoder(r.Body).Decode(&acc)

	err := a.store.UpdateAccount(acc)

	if err != nil{
		return err 
	}

	return WriteJSON(w, http.StatusCreated, acc)
}
 
func (a *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	create := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(create); err != nil {
		return err
	}

	acc, e := types.NewAccount(create.FirstName, create.LastName, create.Password)
	
	if e != nil{
		return e
	}

	if err := a.store.CreateAccount(acc); err != nil {
		return err
	}

	_, err := createJWTToken(acc)

	if err != nil{
		return err 
	}

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

		if tokenStr == ""{
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return
		}
		
		tk, err := validateJWT(tokenStr)

		if !tk.Valid{
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return 
		}

		if err != nil{
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission Denied"})
			return 
		}

		handleFunc(w, r)
	}
}

func getID(r *http.Request) string{
	return mux.Vars(r)["id"]
}

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
