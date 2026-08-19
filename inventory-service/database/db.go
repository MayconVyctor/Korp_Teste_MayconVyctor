package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=admin password=adminpassword dbname=korp_db sslmode=disable"

	var err error

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to configure the database connection:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database did not respond, Check if the Docker container is running: ", err)
	}

	fmt.Println("Connection to PostgreSQL established successfully!")
}

func CreateTable() {

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		code VARCHAR(50) PRIMARY KEY,
		description VARCHAR(255) NOT NULL,
		balance INT NOT NULL
	);`

	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Table created successfully!")
}
