package database_connexion

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Gather PostgreSQL database info
const (
// //host = "localhost"
// host = "/var/run/postgresql"
// //port = 5432
// user = "kwizz"
// // password = ""
// dbname = "kwizz"
)

func Connect() *sql.DB {
	//Crate the connection string
	psqlInfo := "host=/var/run/postgresql dbname=kwizz"
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"dbname=%s sslmode=disable",
	//	host, port, user, dbname)

	//Open a connection to database
	// Open validates the arguments provided but does not create a connection
	// Ping creates the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Successful??
	fmt.Println("Successfully connected!")

	return db
}
