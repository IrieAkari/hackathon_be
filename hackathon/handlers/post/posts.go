package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func PostsGetHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.DB.Query("SELECT id, user_id, content FROM posts WHERE parent_id IS NULL")
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Content); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
