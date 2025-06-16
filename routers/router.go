package routers

import (
    "github.com/kariyana/kakeibo-api/controllers"
    "github.com/kariyana/kakeibo-api/middlewares"
    "net/http"
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
	// 動作確認用の簡単なAPI
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // JWT認証が必要なルート
    auth := router.Group("/")
    auth.Use(middlewares.AuthJWT())
    {
        auth.POST("/kakeibo", controllers.CreateKakeibo)
        auth.GET("/kakeibo", controllers.GetKakeibos)
    }

    router.LoadHTMLGlob("templates/*")
    router.GET("/test-api", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
      "title": "Main website",
    })
  })

    return router
}
