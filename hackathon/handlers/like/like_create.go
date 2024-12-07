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

type CreateLikeRequest struct {
	Email  string `json:"email"`
	PostId string `json:"post_id"`
}

func LikeCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.PostId == "" {
		log.Println("Email or Post ID is empty")
		http.Error(w, "Email or Post ID is empty", http.StatusBadRequest)
		return
	}

	var userId string
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userId)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	likeId := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO likes (id, user_id, post_id) VALUES (?, ?, ?)", likeId, userId, req.PostId)
	if err != nil {
		tx.Rollback()
		log.Printf("Insert error: %v", err)
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE posts SET likes_count = likes_count + 1 WHERE id = ?", req.PostId)
	if err != nil {
		tx.Rollback()
		log.Printf("Update error: %v", err)
		http.Error(w, "Failed to update likes count", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit error: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Like created successfully"))
}
