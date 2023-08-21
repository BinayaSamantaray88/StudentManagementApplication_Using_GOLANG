package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/m/api"
	"example.com/m/db"
	"example.com/m/db2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	createTable(db)

	db2, err := db2.ConnectDB2()
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()
	createTable2(db2)

	r := mux.NewRouter()
	protectedRouter := r.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(api.AuthMiddleware(db2))
	r.HandleFunc("/UserLogin", api.RegisterUser(db)).Methods("POST")
	protectedRouter.HandleFunc("/students", api.CreateStudent(db)).Methods("POST")
	protectedRouter.HandleFunc("/students/{id}", api.GetStudent(db)).Methods("GET")
	protectedRouter.HandleFunc("/students", api.GetALLStudents(db)).Methods("GET")
	protectedRouter.HandleFunc("/students/{id}", api.DeleteStudent(db)).Methods("DELETE")
	protectedRouter.HandleFunc("/students/{id}", api.UpdateStudent(db)).Methods("PUT")
	log.Println("Server started on:8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}

func createTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS student(
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		class INT 
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
func createTable2(db2 *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS user(
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255),
		password VARCHAR(255)
	);`
	_, err := db2.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
