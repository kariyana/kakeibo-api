package main
import (
    "fmt"

    "github.com/kariyana/kakeibo-api/config"
    "github.com/kariyana/kakeibo-api/models"
    "github.com/kariyana/kakeibo-api/routers"

    "golang.org/x/crypto/bcrypt"
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
            Name:     "テストユーザー1",
            Email:    "test1@example.com",
            Password: string(hashed1),
        })
        // テストユーザー2
        pass2 := "password2"
        hashed2, _ := bcrypt.GenerateFromPassword([]byte(pass2), bcrypt.DefaultCost)
        db.Create(&models.User{
            Name:     "テストユーザー2",
            Email:    "test2@example.com",
            Password: string(hashed2),
        })
    }

    // ルーター設定・起動
    r := routers.SetupRouter()
    r.Run("0.0.0.0:8080")
}
