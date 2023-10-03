package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// Helper functions
func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := f(w,r); err!=nil{
			WriteJSON(w, http.StatusBadRequest, APIError{Error:err.Error()})
		}
	}
} 
func getID(r *http.Request) (int, error){
	idStr := mux.Vars(r)["id"]
	id,err := strconv.Atoi(idStr)

	if err!=nil{
		return id, fmt.Errorf("invalid id give %s", idStr)
	}
	return id, nil
}
func WriteJSON(w http.ResponseWriter, status int, v any) error{
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{Error: "permission denied"})
}

// Helper functions for Authorisation through JWT token
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountbyID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}
func validateJWT(tokenString string) (*jwt.Token, error){
	secret:=os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}
func createJWT(account *Account)(string, error){
	claims := &jwt.MapClaims{
		"expiresAt" :15000,
		"accountNumber":account.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}