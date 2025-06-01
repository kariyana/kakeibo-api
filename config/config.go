package config

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
}

var (
	// Cfg はアプリケーション全体で使用する設定
	Cfg Config
)

// LoadEnv は .env から環境変数を読み込む
func LoadEnv() {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println(".env ファイルが見つかりません。環境変数を使用します。")
	// }

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("⚠️ JWT_SECRET が設定されていません")
	}

	Cfg = Config{
		JWTSecret: secret,
	}
}

// Claims はJWTのカスタムクレーム
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
