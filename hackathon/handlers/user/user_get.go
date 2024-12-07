package handlers

import (
	"database/sql"
	"encoding/json"
	"hackathon/models"
	"hackathon/utils"
	"log"
	"net/http"
)

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータを取得
	email := r.URL.Query().Get("email")
	name := r.URL.Query().Get("name")
	id := r.URL.Query().Get("id")

	var rows *sql.Rows
	var err error

	// クエリパラメータに応じて適切なSQLクエリを実行
	switch {
	case email != "":
		rows, err = utils.DB.Query("SELECT id, name, email FROM users WHERE email = ?", email)
	case name != "":
		rows, err = utils.DB.Query("SELECT id, name, email FROM users WHERE name = ?", name)
	case id != "":
		rows, err = utils.DB.Query("SELECT id, name, email FROM users WHERE id = ?", id)
	default:
		// クエリパラメータが指定されていない場合は400 Bad Requestを返す
		log.Println("No query parameter provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// クエリ実行時にエラーが発生した場合の処理
	if err != nil {
		log.Printf("Query error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// クエリ結果を格納するスライス
	var users []models.UserResForHTTPGet
	for rows.Next() {
		var user models.UserResForHTTPGet
		// クエリ結果を構造体にスキャン
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			log.Printf("Scan error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// レスポンスヘッダーを設定し、結果をJSON形式で返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
