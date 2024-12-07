package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

//ユーザーの削除方法を考え直す
//
//
//

func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	// ユーザーの投稿へのいいねを削除
	_, err = tx.Exec("DELETE FROM likes WHERE post_id IN (SELECT id FROM posts WHERE user_id = ?)", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete likes on posts error: %v", err)
		http.Error(w, "Failed to delete likes on posts", http.StatusInternalServerError)
		return
	}

	// ユーザーの投稿を削除
	_, err = tx.Exec("DELETE FROM posts WHERE user_id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete posts error: %v", err)
		http.Error(w, "Failed to delete posts", http.StatusInternalServerError)
		return
	}

	// ユーザーのいいねを削除
	_, err = tx.Exec("DELETE FROM likes WHERE user_id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete likes error: %v", err)
		http.Error(w, "Failed to delete likes", http.StatusInternalServerError)
		return
	}

	// ユーザーを削除
	_, err = tx.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete user error: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit error: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User and related data deleted successfully"))
}
