package handlers

import (
	"database/sql"
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func PostsGetHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	var rows *sql.Rows
	var err error

	if email != "" {
		var userId string
		err = utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userId)
		if err != nil {
			log.Printf("User not found: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		rows, err = utils.DB.Query(`
            SELECT posts.id, posts.user_id, users.name AS user_name, posts.content, posts.likes_count, posts.replys_count, posts.created_at, posts.parent_id
            FROM posts
            JOIN users ON posts.user_id = users.id
            WHERE posts.user_id = ?
            ORDER BY posts.created_at DESC
        `, userId)
	} else {
		rows, err = utils.DB.Query(`
            SELECT posts.id, posts.user_id, users.name AS user_name, posts.content, posts.likes_count, posts.replys_count, posts.created_at, posts.parent_id
            FROM posts
            JOIN users ON posts.user_id = users.id
            WHERE posts.parent_id IS NULL
            ORDER BY posts.created_at DESC
        `)
	}

	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.PostWithUserName
	for rows.Next() {
		var post models.PostWithUserName
		if err := rows.Scan(&post.Id, &post.UserId, &post.UserName, &post.Content, &post.LikesCount, &post.ReplysCount, &post.CreatedAt, &post.ParentId); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
