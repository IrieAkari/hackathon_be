package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "asia-northeast1"      // モデルのリージョン
	modelName = "gemini-1.5-flash-002" // 使用するGeminiモデル名
	projectID = "term6-akari-irie"     // GCPプロジェクトID
)

func TrustScoreReason(content string) (int, string) {
	score, reason, err := getTrustScoreAndDescription(content)
	if err != nil {
		score = -1
		reason = fmt.Sprintf("Error getting trust score and description: %v", err)
	}
	return score, reason
}

func getTrustScoreAndDescription(content string) (int, string, error) {
	// Geminiクライアントの作成
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return -1, "", fmt.Errorf("error creating client: %v", err)
	}

	// プロンプトを作成（信頼度スコアの取得）
	prompt := genai.Text(fmt.Sprintf(
		`%s。この文章の信頼度を0から100の間で評価してください。0は完全に不正確、100は完全に正確を意味します。
		信頼度が判定できない場合でも、最も近い数値を返してください。
		事実や情報を含んでおらずどうしても信頼度が判定できない場合は-1を返してください。
		そしてその信頼度となった理由をその理由を150文字以内で説明してください。
		必ず数値と理由を:で繋げて、数値と理由だけを返してください。何かエラーがあった場合でも、-1:エラーとだけ返してください。`,
		content))

	// Gemini APIを呼び出す
	resp, err := client.GenerativeModel(modelName).GenerateContent(ctx, prompt)
	if err != nil {
		return -1, "", fmt.Errorf("error generating content: %w", err)
	}

	// レスポンスの確認
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return -1, "", fmt.Errorf("no content generated")
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return -1, "", fmt.Errorf("error marshalling resp to JSON: %w", err)
	}

	// JSON から "Parts" の値を取得する
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return -1, "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// "Parts" フィールドの取り出し
	parts, ok := result["Candidates"].([]interface{})[0].(map[string]interface{})["Content"].(map[string]interface{})["Parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return -1, "", fmt.Errorf("invalid or empty parts field")
	}

	// Parts[0] の文字列を取り出し
	output := parts[0].(string)
	//log.Printf("出力結果: %s", output)
	//log.Printf("outputの型 %s", reflect.TypeOf(output))

	// 数値と理由を分割
	partsSplit := strings.SplitN(output, ":", 2)
	if len(partsSplit) != 2 {
		return -1, "", fmt.Errorf("invalid format: expected a score and a reason separated by ':'")
	}

	// 数値部分を score に変換
	score, err := strconv.Atoi(strings.TrimSpace(partsSplit[0]))
	if err != nil {
		return -1, "", fmt.Errorf("error converting score to integer: %w", err)
	}

	// 理由部分を reason に格納
	reason := strings.TrimSpace(partsSplit[1])

	return score, reason, nil
}
