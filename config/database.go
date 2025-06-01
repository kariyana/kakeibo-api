package config

import (
    "fmt"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"

	"time"
)

var DB *gorm.DB

// ConnectDB は MySQL に接続し、gorm.DB インスタンスを返します。
// 環境変数: DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME を使用します。
func ConnectDB() *gorm.DB {
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASS")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    name := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        user, pass, host, port, name)

    var db *gorm.DB
    var err error

    for i := 0; i < 10; i++ {
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err == nil {
            fmt.Println("✅ DB接続成功")
            break
        }
        fmt.Printf("❌ DB接続失敗（%d回目）: %v\n", i+1, err)
        time.Sleep(3 * time.Second)
    }

    if err != nil {
        panic("DB接続エラー: " + err.Error())
    }

    DB = db
    return db
}
