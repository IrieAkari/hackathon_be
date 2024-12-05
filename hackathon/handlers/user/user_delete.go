package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("Email is empty")
		http.Error(w, "Email is empty", http.StatusBadRequest)
		return
	}

	_, err := utils.DB.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		log.Printf("Delete error: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
