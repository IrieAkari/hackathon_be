package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

func LikeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postid")
	email := r.URL.Query().Get("email")
	if postId == "" || email == "" {
		log.Println("Post ID or Email is empty")
		http.Error(w, "Post ID or Email is empty", http.StatusBadRequest)
		return
	}

	var userId string
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userId)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var likeId string
	err = utils.DB.QueryRow("SELECT id FROM likes WHERE post_id = ? AND user_id = ?", postId, userId).Scan(&likeId)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch like ID", http.StatusInternalServerError)
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
