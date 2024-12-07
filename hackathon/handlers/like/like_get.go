package handlers

import (
	"encoding/json"
	"hackathon/utils"
	"log"
	"net/http"
)

type Like struct {
	PostId string `json:"post_id"`
}

func LikeGetHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Println("Email is empty")
		http.Error(w, "Email is empty", http.StatusBadRequest)
		return
	}

	var userId string
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userId)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	rows, err := utils.DB.Query("SELECT post_id FROM likes WHERE user_id = ?", userId)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var likes []Like
	for rows.Next() {
		var like Like
		if err := rows.Scan(&like.PostId); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan like", http.StatusInternalServerError)
			return
		}
		likes = append(likes, like)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(likes)
}
