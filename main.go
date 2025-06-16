package main
import (
    "fmt"

    "github.com/kariyana/kakeibo-api/config"
    "github.com/kariyana/kakeibo-api/models"
    "github.com/kariyana/kakeibo-api/routers"

    "golang.org/x/crypto/bcrypt"
    "math/rand"
	"time"
)

func main() {
    // 環境変数を読み込む
    config.LoadEnv()

    // データベース接続
    db := config.ConnectDB()

    // モデルのマイグレーション（テーブル自動生成）
    db.AutoMigrate(&models.User{}, &models.Kakeibo{})

    // 初期テストユーザー登録（存在しない場合のみ）
    var count int64
    db.Model(&models.User{}).Count(&count)
    if count == 0 {
        fmt.Println("テストユーザーを作成します...")
        // テストユーザー1
        pass1 := "password1"
        hashed1, _ := bcrypt.GenerateFromPassword([]byte(pass1), bcrypt.DefaultCost)
        db.Create(&models.User{
            PublicID: GenerateRandomString(9), // ランダムな公開IDを生成
            UserName:     "テストユーザー1",
            Email:    "test1@example.com",
            PasswordHash: string(hashed1),
        })
        // テストユーザー2
        pass2 := "password2"
        hashed2, _ := bcrypt.GenerateFromPassword([]byte(pass2), bcrypt.DefaultCost)
        db.Create(&models.User{
            PublicID: GenerateRandomString(9), // ランダムな公開IDを生成
            UserName:     "テストユーザー2",
            Email:    "test2@example.com",
            PasswordHash: string(hashed2),
        })
    }

    // ルーター設定・起動
    r := routers.SetupRouter()
    r.Run("0.0.0.0:8080")
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	// 毎回異なる乱数を生成するために seed を初期化
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
