# ================================
# ビルドステージ
# ================================
FROM golang:1.21-alpine AS builder

# 必要なパッケージをインストール（git, ca-certificates）
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 依存関係ファイルを先にコピー（キャッシュ最適化）
COPY go.mod ./
COPY go.sum ./

# main.go だけ先にコピー（go mod tidy に必要な import を含む）
# または全体の .go ファイル群の中の import に依存するものをコピー
COPY main.go ./

# 一部の依存関係を含むディレクトリも先にコピー（必要に応じて）
COPY config/ ./config/
COPY controllers/ ./controllers/
COPY routers/ ./routers/
COPY middlewares/ ./middlewares/
COPY models/ ./models/

# 依存関係を解決し、go.sum を整える
RUN go mod tidy

# アプリケーションソースをコピー
COPY . .

# 静的リンクでビルド（軽量ランタイムと互換性あり）
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# ================================
# 実行ステージ
# ================================
FROM alpine:latest

# ca-certificates（TLSのため）をインストール
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# builderからビルド済みバイナリをコピー
COPY --from=builder /app/app .

# バイナリを起動（ポートは docker-compose.yml 側で指定）
CMD ["./app"]
