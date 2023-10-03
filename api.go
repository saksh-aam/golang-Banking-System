package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store: store,
	}
}

func (s *APIServer) Run(){
	router :=mux.NewRouter()

	router.HandleFunc("/signup", makeHTTPHandleFunc(s.handleSignup))
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountbyID), s.store))
	router.HandleFunc("/transfer/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleTransfer), s.store))
	router.HandleFunc("/addFund/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleAddFunds), s.store))
	log.Println("Server listening on port", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}



