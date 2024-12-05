package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

func ReplyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	replyId := r.URL.Query().Get("replyid")
	if replyId == "" {
		log.Println("Reply ID is empty")
		http.Error(w, "Reply ID is empty", http.StatusBadRequest)
		return
	}

	_, err := utils.DB.Exec("DELETE FROM posts WHERE id = ?", replyId)
	if err != nil {
		log.Printf("Delete error: %v", err)
		http.Error(w, "Failed to delete reply", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reply deleted successfully"))
}
