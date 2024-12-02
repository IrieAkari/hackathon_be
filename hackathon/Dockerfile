# ベースイメージとして Go 1.23 を使用
FROM golang:1.23

# 作業ディレクトリを設定
WORKDIR /app

# Go モジュールの初期化と依存関係のダウンロード
COPY go.mod ./
# go.sum ファイルが存在しない場合の対処
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# アプリケーションをビルド
RUN go build -o main .



# ポート 8080 を公開
EXPOSE 8080

# コンテナ起動時に実行するコマンド
CMD ["./main"]