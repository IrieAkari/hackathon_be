package main

import (
	"log"
	"net/http"
	"os"
	_ "os/signal"
	_ "syscall"

	"hackathon/handlers"
	"hackathon/utils"

	"github.com/rs/cors"
)

func main() {
	utils.InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/user", handlers.UserGetHandler)
	mux.HandleFunc("/register", handlers.RegisterPostHandler)
	mux.HandleFunc("/users", handlers.UsersGetHandler)

	handler := cors.Default().Handler(mux)

	utils.CloseDBWithSysCall()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
