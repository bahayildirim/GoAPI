package router

import (
	"mods/api"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) { // Registers routes and CRUD operations to URL's
	router.HandleFunc("/users", api.GetUsers).Methods("GET")
	router.HandleFunc("/user", api.PostUser).Methods("POST")
	router.HandleFunc("/user/{id}", api.GetUser_Id).Methods("GET")
	router.HandleFunc("/user/{id}", api.DeleteUser).Methods("DELETE")
}
