package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func UserEmailGetHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("Email is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is empty"))
		return
	}

	rows, err := utils.DB.Query("SELECT id, name, email FROM users WHERE email = ?", email)
	if err != nil {
		log.Printf("Query error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Query error"))
		return
	}
	defer rows.Close()

	var users []models.UserResForHTTPGet
	for rows.Next() {
		var user models.UserResForHTTPGet
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			log.Printf("Scan error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Scan error"))
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
