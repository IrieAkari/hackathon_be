package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

//リプライも削除する
//再帰
//
//
//
//
//

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postid")
	if postId == "" {
		log.Println("Post ID is empty")
		http.Error(w, "Post ID is empty", http.StatusBadRequest)
		return
	}

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	// 投稿へのいいねを削除
	_, err = tx.Exec("DELETE FROM likes WHERE post_id = ?", postId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete likes error: %v", err)
		http.Error(w, "Failed to delete likes", http.StatusInternalServerError)
		return
	}

	// ポストを削除
	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", postId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete post error: %v", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Transaction commit error: %v", err)
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post and related likes deleted successfully"))
}
