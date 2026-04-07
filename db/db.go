package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

//go:embed schema/base.sql
var schema string

func connectDB() (*sql.DB, error) {
	fmt.Println("Connecting to database...")

	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var user = os.Getenv("POSTGRES_USER")
	var password = os.Getenv("POSTGRES_PASSWORD")
	var dbname = os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db := waitForDB(psqlInfo)
	if db == nil {
		return nil, fmt.Errorf("could not connect to database")
	}

	fmt.Println("connected to db")

	return db, nil
}

func waitForDB(dbURL string) *sql.DB {
	var db *sql.DB
	var err error
	for range 10 {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			if err = db.Ping(); err == nil {
				return db
			}
		}
		log.Println("waiting for db...")
		time.Sleep(2 * time.Second)
	}
	log.Fatal("could not connect to db")
	return nil
}

func loadSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

func InitDB() (*sql.DB, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

	err = loadSchema(db)
	if err != nil {
		return nil, err
	}

	fmt.Println("schema loaded")

	return db, nil
}
