package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := mux.NewRouter()
	http.Handle("/", router)
	registerRoutes(router)
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) { // Display all users
	w.WriteHeader(http.StatusOK)
	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	fmt.Printf("got /users request\n")
	row, err := db.Query("SELECT * FROM users ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Cycle through entries on database
		var id int
		var name string
		var surname string
		var age int
		row.Scan(&id, &name, &surname, &age)
		fmt.Fprintf(w, name+" "+surname+" "+strconv.Itoa(age)+"\n")
	}
}

func getUser_Id(w http.ResponseWriter, r *http.Request) { // Display users with given ID
	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	tempId := mux.Vars(r)["id"]
	dbId, _ := strconv.Atoi(tempId)
	fmt.Printf("got /user/%s request\n", tempId)
	row := db.QueryRow("SELECT * FROM users WHERE (id = ?)", dbId)
	var id int
	var name string
	var surname string
	var age int
	row.Scan(&id, &name, &surname, &age)
	io.WriteString(w, name+" "+surname+" "+strconv.Itoa(age)+"\n")
}

func postUser(w http.ResponseWriter, r *http.Request) { // Post a user inputted on user.html
	name := "Baha"
	surname := "Yıldırım"
	age := 21
	fmt.Fprintf(w, "POST request successful\n")
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Surname = %s\n", surname)
	fmt.Fprintf(w, "Age = %d\n", age)
	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	insertUserSQL := `INSERT INTO users(name, surname, age) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL) // Prepare statement.
	_, err = statement.Exec(name, surname, age)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	tempId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(tempId)
	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	insertUserSQL := `DELETE FROM users WHERE (id = ?)`
	statement, err := db.Prepare(insertUserSQL) // Prepare statement.
	_, err = statement.Exec(id)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Fprintf(w, "Entry with id %d deleted\n", id)
}

func registerRoutes(router *mux.Router) { // Registers routes and CRUD operations to URL's
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/user", postUser).Methods("POST")
	router.HandleFunc("/user/{id}", getUser_Id).Methods("GET")
	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
}
