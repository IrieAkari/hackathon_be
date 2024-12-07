package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

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

	// ユーザーの投稿IDと親投稿IDを取得
	rows, err := tx.Query("SELECT id, parent_id FROM posts WHERE user_id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Query posts error: %v", err)
		http.Error(w, "Failed to query posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var postIds []string
	var parentIds []string
	for rows.Next() {
		var postId string
		var parentId *string
		if err := rows.Scan(&postId, &parentId); err != nil {
			tx.Rollback()
			log.Printf("Scan post ID error: %v", err)
			http.Error(w, "Failed to scan post ID", http.StatusInternalServerError)
			return
		}
		postIds = append(postIds, postId)
		if parentId != nil {
			parentIds = append(parentIds, *parentId)
		}
	}

	// 各投稿について、その投稿を親投稿とする投稿のis_parent_deletedをTrueに設定
	for _, postId := range postIds {
		_, err = tx.Exec("UPDATE posts SET is_parent_deleted = TRUE WHERE parent_id = ?", postId)
		if err != nil {
			tx.Rollback()
			log.Printf("Update is_parent_deleted error: %v", err)
			http.Error(w, "Failed to update is_parent_deleted", http.StatusInternalServerError)
			return
		}
	}

	// ユーザーのいいねを削除し、その対象の投稿のいいね数を減らす
	likeRows, err := tx.Query("SELECT post_id FROM likes WHERE user_id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Query likes error: %v", err)
		http.Error(w, "Failed to query likes", http.StatusInternalServerError)
		return
	}
	defer likeRows.Close()

	var likedPostIds []string
	for likeRows.Next() {
		var postId string
		if err := likeRows.Scan(&postId); err != nil {
			tx.Rollback()
			log.Printf("Scan like post ID error: %v", err)
			http.Error(w, "Failed to scan like post ID", http.StatusInternalServerError)
			return
		}
		likedPostIds = append(likedPostIds, postId)
	}

	for _, postId := range likedPostIds {
		_, err = tx.Exec("UPDATE posts SET likes_count = likes_count - 1 WHERE id = ?", postId)
		if err != nil {
			tx.Rollback()
			log.Printf("Update likes_count error: %v", err)
			http.Error(w, "Failed to update likes_count", http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.Exec("DELETE FROM likes WHERE user_id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete likes error: %v", err)
		http.Error(w, "Failed to delete likes", http.StatusInternalServerError)
		return
	}

	// ユーザーの投稿へのいいねを削除
	for _, postId := range postIds {
		_, err = tx.Exec("DELETE FROM likes WHERE post_id = ?", postId)
		if err != nil {
			tx.Rollback()
			log.Printf("Delete likes on posts error: %v", err)
			http.Error(w, "Failed to delete likes on posts", http.StatusInternalServerError)
			return
		}
	}

	// ユーザーの投稿を削除
	_, err = tx.Exec("DELETE FROM posts WHERE user_id = ?", userId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete posts error: %v", err)
		http.Error(w, "Failed to delete posts", http.StatusInternalServerError)
		return
	}

	// 親投稿のリプライ数を-1
	for _, parentId := range parentIds {
		_, err = tx.Exec("UPDATE posts SET replys_count = replys_count - 1 WHERE id = ?", parentId)
		if err != nil {
			tx.Rollback()
			log.Printf("Update replys_count error: %v", err)
			http.Error(w, "Failed to update replys_count", http.StatusInternalServerError)
			return
		}
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
