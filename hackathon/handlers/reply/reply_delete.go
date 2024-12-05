package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

//リプライへのリプライも削除
//
//
//
//
//

func ReplyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	replyId := r.URL.Query().Get("replyid")
	if replyId == "" {
		log.Println("Reply ID is empty")
		http.Error(w, "Reply ID is empty", http.StatusBadRequest)
		return
	}

	tx, err := utils.DB.Begin()
	if err != nil {
		log.Printf("Transaction begin error: %v", err)
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	// 親投稿のIDを取得
	var parentId string
	err = tx.QueryRow("SELECT parent_id FROM posts WHERE id = ?", replyId).Scan(&parentId)
	if err != nil {
		tx.Rollback()
		log.Printf("Query parent_id error: %v", err)
		http.Error(w, "Failed to get parent_id", http.StatusInternalServerError)
		return
	}

	// リプライへのいいねを削除
	_, err = tx.Exec("DELETE FROM likes WHERE post_id = ?", replyId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete likes error: %v", err)
		http.Error(w, "Failed to delete likes", http.StatusInternalServerError)
		return
	}

	// リプライを削除
	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", replyId)
	if err != nil {
		tx.Rollback()
		log.Printf("Delete reply error: %v", err)
		http.Error(w, "Failed to delete reply", http.StatusInternalServerError)
		return
	}

	// 親投稿のリプライ数を-1
	_, err = tx.Exec("UPDATE posts SET replys_count = replys_count - 1 WHERE id = ?", parentId)
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
	w.Write([]byte("Reply and related likes deleted successfully"))
}
