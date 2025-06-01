# バイナリをビルド
build:
	go build -o bin/app main.go

# ローカルで実行
run:
	go run main.go

# Dockerイメージのビルドと起動
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# モジュールの整理
tidy:
	go mod tidy
