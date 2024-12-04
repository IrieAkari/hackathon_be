package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables...")
	}

	dbUser := os.Getenv("MYSQL_USER")
	dbPwd := os.Getenv("MYSQL_PWD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbHost := os.Getenv("MYSQL_HOST")

	dbURI := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", dbUser, dbPwd, dbHost, dbName)
	log.Printf("Connecting to database with URI: %s", dbURI)

	var errDB error
	DB, errDB = sql.Open("mysql", dbURI)
	if errDB != nil {
		log.Fatalf("Failed to open database: %v", errDB)
	}

	if errPing := DB.Ping(); errPing != nil {
		log.Println("Failed to ping database: %v", errPing)
	}

	log.Println("Connected to database successfully")
}

func CloseDBWithSysCall() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Closing database connection")
		if err := DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		os.Exit(0)
	}()
}
