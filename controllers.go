package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//  (s *APIServer is a method receiver, which allows us to attach a function to a struct or any other custom type)
func (s *APIServer) handleSignup(w http.ResponseWriter, r *http.Request) error {
	if(r.Method !="POST"){
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	var req SignUpRequest
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil{
		return err
	}
	acc, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account => ", acc.Number)
	return WriteJSON(w, http.StatusOK, acc)
}
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountbyNumber(int(req.Number))
	if err != nil {
		return err
	}

	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	return WriteJSON(w, http.StatusOK, resp)
}
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method =="GET"{
		return s.handleGetAccount(w,r)
	}
	if r.Method =="POST"{
		return s.handleCreateAccount(w,r)
	}
	return nil
}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err!=nil{
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *APIServer) handleGetAccountbyID(w http.ResponseWriter, r *http.Request) error {
	if(r.Method=="GET"){
		id, err :=getID(r)

		if err!=nil{
			return err
		}
		account, err := s.store.GetAccountbyID(id)
		if err!=nil{
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method =="DELETE"{
		return s.handleDeleteAccount(w,r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err:= json.NewDecoder(r.Body).Decode(createAccountReq); err!=nil{
		return err
	}
	account, err :=NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Password)
	if err!=nil{
		return err
	}
	if err:=s.store.CreateAccount(account); err!=nil{
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err :=getID(r)

	if err!=nil{
		return err
	}

	if err :=s.store.DeleteAccount(id); err!=nil{
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string] int{"deleted":id} )
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	id, err :=getID(r)
	if err!=nil{
		return err
	}

	sender, err := s.store.GetAccountbyID(id)
	if err!=nil{
		return fmt.Errorf("sender account id %v doesn't exist in db", id)
	}
	receiver, err := s.store.GetAccountbyNumber(transferReq.ToAccount)
	if err!=nil{
		return fmt.Errorf("receiver account number %v doesn't exist in db", transferReq.ToAccount)
	}
	if(sender.Balance < int64(transferReq.Amount)){
		return fmt.Errorf("sender doesn't have much balance to transfer")
	}

	s.store.UpdateAccount(int(sender.Number), int(sender.Balance) - transferReq.Amount)
	s.store.UpdateAccount(int(receiver.Number), int(receiver.Balance) + transferReq.Amount)

	var resp TransferReponse
	json.Unmarshal([]byte(`{"Message":"Amount successfully transfered!"}`), &resp)
	return WriteJSON(w, http.StatusOK, resp)
}
func (s *APIServer) handleAddFunds(w http.ResponseWriter, r *http.Request) error {
	addFundReq := new(AddFundRequest)
	if err := json.NewDecoder(r.Body).Decode(addFundReq); err != nil {
		return err
	}

	defer r.Body.Close()
	id, err :=getID(r)
	if err!=nil{
		return err
	}

	account, err := s.store.GetAccountbyID(id)

	if err!=nil{
		return fmt.Errorf("account id %v doesn't exist in db", id)
	}

	s.store.UpdateAccount(int(account.Number), int(account.Balance) + addFundReq.Amount)

	var resp TransferReponse
	json.Unmarshal([]byte(`{"message":"Amount added successfully to your account!"}`), &resp)
	return WriteJSON(w, http.StatusOK, resp)
}