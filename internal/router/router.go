package router

import (
	"ROOmail/internal/handlers/auth"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")

	return r
}
