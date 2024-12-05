package handlers

import (
	"hackathon/utils"
	"log"
	"net/http"
)

func PostDeleteHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postid")
	if postId == "" {
		log.Println("Post ID is empty")
		http.Error(w, "Post ID is empty", http.StatusBadRequest)
		return
	}

	_, err := utils.DB.Exec("DELETE FROM posts WHERE id = ?", postId)
	if err != nil {
		log.Printf("Delete error: %v", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}
