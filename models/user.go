package models

import (
    "gorm.io/gorm"
)

// User はユーザー情報を表すモデルです。
type User struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"-"`
}
