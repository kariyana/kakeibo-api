package middlewares

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"

    "github.com/kariyana/kakeibo-api/config"
)

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
            c.Abort()
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(config.Cfg.JWTSecret), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
            c.Abort()
            return
        }
        claims := token.Claims.(jwt.MapClaims)
        userID := uint(claims["user_id"].(float64))
        c.Set("user_id", userID)
        c.Next()
    }
}
