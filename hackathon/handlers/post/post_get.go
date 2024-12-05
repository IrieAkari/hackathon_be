package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func PostGetHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postid")
	if postId == "" {
		log.Println("Post ID is empty")
		http.Error(w, "Post ID is empty", http.StatusBadRequest)
		return
	}

	var post models.PostWithUserName
	err := utils.DB.QueryRow(`
        SELECT posts.id, posts.user_id, users.name AS user_name, posts.content, posts.created_at, posts.parent_id
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.id = ?
    `, postId).Scan(&post.Id, &post.UserId, &post.UserName, &post.Content, &post.CreatedAt, &post.ParentId)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
