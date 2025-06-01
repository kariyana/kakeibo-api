package routers

import (
    "kakeibo/controllers"
    "kakeibo/middlewares"

    "github.com/gin-gonic/gin"
)

// SetupRouter はGinのルーティングを設定してエンジンを返します。
func SetupRouter() *gin.Engine {
    router := gin.Default()

    // 認証不要なルート
    router.POST("/signup", controllers.Signup)
    router.POST("/login", controllers.Login)
    router.GET("/auth/google", controllers.GoogleLogin)
    router.GET("/auth/google/callback", controllers.GoogleCallback)

    // JWT認証が必要なルート
    auth := router.Group("/")
    auth.Use(middlewares.AuthJWT())
    {
        auth.POST("/kakeibo", controllers.CreateKakeibo)
        auth.GET("/kakeibo", controllers.GetKakeibos)
    }

    return router
}
