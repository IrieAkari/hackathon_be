package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		log.Println("Name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query("SELECT id, name, email FROM users WHERE name = ?", name)
	if err != nil {
		log.Printf("Query error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResForHTTPGet
	for rows.Next() {
		var user models.UserResForHTTPGet
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			log.Printf("Scan error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
