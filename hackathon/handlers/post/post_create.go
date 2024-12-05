package handlers

import (
	"encoding/json"
	"hackathon/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

type CreatePostRequest struct {
	Email   string `json:"email"`
	Content string `json:"content"`
}

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Content == "" {
		log.Println("Email or content is empty")
		http.Error(w, "Email or content is empty", http.StatusBadRequest)
		return
	}

	var userID string
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	postID := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

	_, err = utils.DB.Exec("INSERT INTO posts (id, user_id, content) VALUES (?, ?, ?)", postID, userID, req.Content)
	if err != nil {
		log.Printf("Insert error: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post created successfully"))
}
