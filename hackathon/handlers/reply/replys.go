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

	rows, err := utils.DB.Query("SELECT id, user_id, content FROM posts WHERE parent_id = ?", parentId)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch replies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var replies []models.Post
	for rows.Next() {
		var reply models.Post
		if err := rows.Scan(&reply.Id, &reply.UserId, &reply.Content); err != nil {
			log.Printf("Scan error: %v", err)
			http.Error(w, "Failed to scan reply", http.StatusInternalServerError)
			return
		}
		replies = append(replies, reply)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}
