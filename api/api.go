package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Database struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) { // Display all users
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

func GetUser_Id(w http.ResponseWriter, r *http.Request) { // Display users with given ID
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

func PostUser(w http.ResponseWriter, r *http.Request) { // Post a user inputted on user.html
	var database Database
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&database)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	fmt.Fprintf(w, "POST request successful\n")
	fmt.Fprintf(w, "Name = %s\n", database.Name)
	fmt.Fprintf(w, "Surname = %s\n", database.Surname)
	fmt.Fprintf(w, "Age = %d\n", database.Age)
	db, _ := sql.Open("sqlite3", "./sqlite-database.db")
	insertUserSQL := `INSERT INTO users(name, surname, age) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL) // Prepare statement.
	_, err = statement.Exec(database.Name, database.Surname, database.Age)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
