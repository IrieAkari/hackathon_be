package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func PostsGetHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.DB.Query(`
        SELECT posts.id, posts.user_id, users.name AS user_name, posts.content 
        FROM posts 
        JOIN users ON posts.user_id = users.id 
        WHERE posts.parent_id IS NULL
    `)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.PostWithUserName
	for rows.Next() {
		var post models.PostWithUserName
		if err := rows.Scan(&post.Id, &post.UserId, &post.UserName, &post.Content); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
