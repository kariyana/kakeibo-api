package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement;comment:ID"`
	PublicID     string         `gorm:"size:9;not null;unique;comment:公開ID"`
	UserName     string         `gorm:"size:20;not null;comment:ユーザー名"`
	Email        string         `gorm:"size:255;not null;unique;comment:メールアドレス"`
	PasswordHash string         `gorm:"size:255;not null;comment:パスワード"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:論理削除"`
	CreatedAt    time.Time      `gorm:"not null;comment:作成日時"`
	UpdatedAt    time.Time      `gorm:"not null;comment:更新日時"`
}
