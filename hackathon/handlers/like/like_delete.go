package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

func LikeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	likeId := r.URL.Query().Get("likeid")
	if likeId == "" {
		log.Println("Like ID is empty")
		http.Error(w, "Like ID is empty", http.StatusBadRequest)
		return
	}

	var postId string
	err := utils.DB.QueryRow("SELECT post_id FROM likes WHERE id = ?", likeId).Scan(&postId)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch post ID", http.StatusInternalServerError)
		return
	}

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM likes WHERE id = ?", likeId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete error: %v", err)
		http.Error(w, "Failed to delete like", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE posts SET likes_count = likes_count - 1 WHERE id = ?", postId)
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
	w.Write([]byte("Like deleted successfully"))
}
