package models

import (
    "gorm.io/gorm"
)

// Kakeibo は家計簿のエントリー（収支項目）を表すモデルです。
type Kakeibo struct {
    gorm.Model
    UserID uint   `json:"user_id"`
    Item   string `json:"item"`
    Amount int    `json:"amount"`
    Memo   string `json:"memo"`
}
