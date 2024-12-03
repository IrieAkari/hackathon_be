package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid"
	"github.com/rs/cors"
)

type UserResForHTTPGet struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserReqForHTTPPost struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var DB *sql.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		if name == "" {
			log.Println("Name is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rows, err := DB.Query("SELECT id, name, age FROM user WHERE name = ?", name)
		if err != nil {
			log.Printf("Query error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []UserResForHTTPGet
		for rows.Next() {
			var user UserResForHTTPGet
			if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
				log.Printf("Scan error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var user UserReqForHTTPPost
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Decode error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

		_, err := DB.Exec("INSERT INTO user (id, name, age) VALUES (?, ?, ?)", id, user.Name, user.Age)
		if err != nil {
			log.Printf("Insert error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": id})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func closeDBWithSysCall() {
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

func main() {

	initDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/user", handler)

	handler := cors.Default().Handler(mux)

	closeDBWithSysCall()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
