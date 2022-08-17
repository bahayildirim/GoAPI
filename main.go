package main

import (
	"mods/db"
	"mods/router"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	rt := mux.NewRouter()
	router.RegisterRoutes(rt)
	db.ConnectDatabase(rt)
}
