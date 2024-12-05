package handlers

import (
	"database/sql"
	"encoding/json"
	"hackathon/utils"
	"log"
	"net/http"
)

type GetLikeRequest struct {
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
}

type GetLikeResponse struct {
	LikeId *string `json:"like_id"`
}

func LikeGetHandler(w http.ResponseWriter, r *http.Request) {
	var req GetLikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserId == "" || req.PostId == "" {
		log.Println("User ID or Post ID is empty")
		http.Error(w, "User ID or Post ID is empty", http.StatusBadRequest)
		return
	}

	var likeId *string
	err := utils.DB.QueryRow("SELECT id FROM likes WHERE user_id = ? AND post_id = ?", req.UserId, req.PostId).Scan(&likeId)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch like", http.StatusInternalServerError)
		return
	}

	response := GetLikeResponse{
		LikeId: likeId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
