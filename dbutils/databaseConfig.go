package dbutils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

// default configuration values if environment variables are not set
const (
	defaultHost     = "0.0.0.0"
	defaultPort     = 5432
	defaultUser     = "postgres"
	defaultPassword = "rocket123"
	defaultDBName   = "projectDb"
)

type DB struct {
	*sql.DB
}

func ConnectDatabase() (*DB, error) {
	// read configuration from environment with defaults
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	portStr := os.Getenv("DB_PORT")
	port := defaultPort
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = defaultUser
	}

	password := os.Getenv("DB_PASS")
	if password == "" {
		password = defaultPassword
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = defaultDBName
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	fmt.Println(psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB connected")

	sqlStmt, err := os.ReadFile("./db/init.sql")
	if err != nil {
		log.Fatalf("Error on reading init.sql: %s", err)
	}

	_, err = db.Exec(string(sqlStmt))
	if err != nil {
		log.Fatalf("err exec: %s", err)
	}

	fmt.Println("Database Tables initialized")

	return &DB{db}, nil

}
