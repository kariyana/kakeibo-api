package config

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
)

var JWTSecret string

// LoadEnv は .env ファイルを読み込み、JWT_SECRETなどを環境変数から設定します。
func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("dotenvファイルが見つかりません。環境変数を使用します。")
    }
    JWTSecret = os.Getenv("JWT_SECRET")
    if JWTSecret == "" {
        fmt.Println("Warning: JWT_SECRET が設定されていません。")
    }
}
