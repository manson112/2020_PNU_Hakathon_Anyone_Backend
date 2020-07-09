package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// InitDatabase :: Initializing database
func InitDatabase() {
	log.Println("Initializing Database[Mysql] ...")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	connectionInfo := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	db, err := sql.Open("mysql", connectionInfo)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	log.Println("Database[Mysql] successfully initialized")
}
