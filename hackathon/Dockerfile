#goを動かすためのDockerfile

# Build stage
FROM golang:1.22.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

# Runtime stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main /app/main

#COPY application_default_credentials.json /app/application_default_credentials.json

# 環境変数を設定
#ENV GOOGLE_APPLICATION_CREDENTIALS="/app/application_default_credentials.json"

CMD ["./main"]