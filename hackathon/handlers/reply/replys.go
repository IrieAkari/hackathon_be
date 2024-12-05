package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func ReplysGetHandler(w http.ResponseWriter, r *http.Request) {
	parentId := r.URL.Query().Get("parentid")
	if parentId == "" {
		log.Println("Parent ID is empty")
		http.Error(w, "Parent ID is empty", http.StatusBadRequest)
		return
	}

	rows, err := utils.DB.Query(`
        SELECT posts.id, posts.user_id, users.name AS user_name, posts.content, posts.likes_count, posts.replys_count, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
        WHERE posts.parent_id = ?
        ORDER BY posts.created_at DESC
    `, parentId)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch replies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var replies []models.PostWithUserName
	for rows.Next() {
		var reply models.PostWithUserName
		if err := rows.Scan(&reply.Id, &reply.UserId, &reply.UserName, &reply.Content, &reply.LikesCount, &reply.ReplysCount, &reply.CreatedAt); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan reply", http.StatusInternalServerError)
			return
		}
		replies = append(replies, reply)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}
