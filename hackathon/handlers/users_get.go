package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func UsersGetHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := utils.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Printf("fail: db.Query, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := make([]models.UserResForHTTPGet, 0)
	for rows.Next() {
		var u models.UserResForHTTPGet
		if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
			log.Printf("fail: rows.Scan, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
