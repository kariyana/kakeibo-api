package controllers

import (
    "net/http"
    // "time"

    "github.com/kariyana/kakeibo-api/config"
    "github.com/kariyana/kakeibo-api/models"

    "github.com/gin-gonic/gin"
)

// CreateKakeibo は家計簿エントリーを登録します。認証されたユーザーのみアクセス可。
func CreateKakeibo(c *gin.Context) {
    var input struct {
        Item   string `json:"item"`
        Amount int    `json:"amount"`
        Memo   string `json:"memo"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "パラメータが正しくありません。"})
        return
    }
    // コンテキストからユーザーIDを取得
    userIDVal, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "認証情報がありません。"})
        return
    }
    userID := userIDVal.(uint)

    entry := models.Kakeibo{
        UserID: userID,
        Item:   input.Item,
        Amount: input.Amount,
        Memo:   input.Memo,
    }
    if err := config.DB.Create(&entry).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "エントリーの登録に失敗しました。"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "家計簿エントリーを追加しました。"})
}

// GetKakeibos は認証ユーザーの家計簿エントリー一覧を取得します。
func GetKakeibos(c *gin.Context) {
    userIDVal, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "認証情報がありません。"})
        return
    }
    userID := userIDVal.(uint)

    var entries []models.Kakeibo
    if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&entries).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "エントリーの取得に失敗しました。"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "家計簿エントリーを取得しました。",
        "data":    entries,
    })
}
