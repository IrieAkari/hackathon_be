package handlers

import (
	"encoding/json"
	"hackathon/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"hackathon/handlers/gemini" // TrustScoreReason関数をインポート

	"github.com/oklog/ulid"
)

type CreateReplyRequest struct {
	Email    string `json:"email"`
	Content  string `json:"content"`
	ParentId string `json:"parent_id"`
}

func ReplyCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Content == "" || req.ParentId == "" {
		log.Println("Email, content, or parent_id is empty")
		http.Error(w, "Email, content, or parent_id is empty", http.StatusBadRequest)
		return
	}

	var userId string
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userId)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	replyId := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

	// 信頼度スコアと説明を生成
	trustScore, trustDescription := gemini.TrustScoreReason(req.Content)

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO posts (id, user_id, content, parent_id, trust_score, trust_description) VALUES (?, ?, ?, ?, ?, ?)", replyId, userId, req.Content, req.ParentId, trustScore, trustDescription)
	if err != nil {
		tx.Rollback()
		log.Printf("Insert error: %v", err)
		http.Error(w, "Failed to create reply", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE posts SET replys_count = replys_count + 1 WHERE id = ?", req.ParentId)
	if err != nil {
		tx.Rollback()
		log.Printf("Update replys_count error: %v", err)
		http.Error(w, "Failed to update replys_count", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit error: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reply created successfully"))
}
