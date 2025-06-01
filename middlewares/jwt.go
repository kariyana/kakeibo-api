package middlewares

import (
    "net/http"
    "strings"

    "kakeibo/config"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

// AuthJWT はJWT認証ミドルウェアです。
func AuthJWT() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "認証が必要です。"})
            c.Abort()
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // トークン解析
        token, err := jwt.ParseWithClaims(tokenString, &config.Claims{}, func(t *jwt.Token) (interface{}, error) {
            return []byte(config.JWTSecret), nil
        })
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "不正なトークンです。"})
            c.Abort()
            return
        }
        claims, ok := token.Claims.(*config.Claims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "認証に失敗しました。"})
            c.Abort()
            return
        }
        // コンテキストにユーザーIDを設定
        c.Set("userID", claims.UserID)
        c.Next()
    }
}
