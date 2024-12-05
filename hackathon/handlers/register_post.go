package handlers

import (
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserReqForHTTPPost
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

	_, err := utils.DB.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)", id, user.Name, user.Email)
	if err != nil {
		log.Printf("Insert error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
