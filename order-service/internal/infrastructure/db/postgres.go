package db

import (
	"database/sql"
	"fmt"
	_ "github.com/joho/godotenv"
	"log"
	"os"
)

func ConnectPostgres() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println("Connecting to DB with DSN:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Printf("Connected to PostgreSQL")
	return db, nil
}
