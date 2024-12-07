package handlers

import (
	"encoding/json"
	"hackathon/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"hackathon/handlers/gemini" // TrustScoreReason関数をインポート

	"github.com/oklog/ulid"
)

// クライアントからのリクエストの構造体を定義
type CreatePostRequest struct {
	Email   string `json:"email"`   // 投稿者のメールアドレス
	Content string `json:"content"` // 投稿内容
}

// 投稿作成のHTTPハンドラー
func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest

	// リクエストボディをデコード
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 必須フィールドの検証
	if req.Email == "" || req.Content == "" {
		log.Println("Email or content is empty")
		http.Error(w, "Email or content is empty", http.StatusBadRequest)
		return
	}

	// ユーザーIDを取得
	var userID string
	err := utils.DB.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// ULIDを使用して投稿IDを生成
	postID := ulid.MustNew(ulid.Timestamp(time.Now()), rand.New(rand.NewSource(time.Now().UnixNano()))).String()

	// 信頼度スコアと説明を生成
	trustScore, trustDescription := gemini.TrustScoreReason(req.Content)

	// データベースに投稿を保存
	_, err = utils.DB.Exec("INSERT INTO posts (id, user_id, content, trust_score, trust_description) VALUES (?, ?, ?, ?, ?)", postID, userID, req.Content, trustScore, trustDescription)
	if err != nil {
		log.Printf("Insert error: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// 成功レスポンスを返す
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post created successfully"))
}
