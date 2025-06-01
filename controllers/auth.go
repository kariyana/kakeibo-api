package controllers

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "time"

    "kakeibo/config"
    "kakeibo/models"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/oauth2"
    _ "github.com/joho/godotenv/autoload"
)

// Google OAuth2 設定（.env から CLIENT ID/SECRET/REDIRECT_URL を取得）
var googleOauthConfig = &oauth2.Config{
    RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
    Endpoint:     google.Endpoint,
}

// Claims は JWT のペイロード（主張）を定義します。
type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

// Signup は新規ユーザー登録を行います。
func Signup(c *gin.Context) {
    var input struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "パラメータが正しくありません。"})
        return
    }
    // 同じメールアドレスのユーザーが存在しないかチェック
    var existing models.User
    if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "既に登録されているメールアドレスです。"})
        return
    }
    // パスワードをハッシュ化して保存
    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "サーバーエラーが発生しました。"})
        return
    }
    user := models.User{Name: input.Name, Email: input.Email, Password: string(hashed)}
    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "ユーザー登録に失敗しました。"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "ユーザー登録に成功しました。"})
}

// Login はメールアドレスとパスワードでログインし、JWTを発行します。
func Login(c *gin.Context) {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "パラメータが正しくありません。"})
        return
    }
    var user models.User
    if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "メールアドレスまたはパスワードが違います。"})
        return
    }
    // パスワード照合
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "メールアドレスまたはパスワードが違います。"})
        return
    }
    // JWT トークン生成（有効期限24時間）
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(config.JWTSecret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "トークン生成に失敗しました。"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "ログインに成功しました。",
        "token":   tokenString,
    })
}

// GoogleLogin はGoogle OAuthの認証URLへリダイレクトします。
func GoogleLogin(c *gin.Context) {
    url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback はGoogleからのコールバックを処理し、ユーザーを作成または更新してJWTを返します。
func GoogleCallback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "認証コードが取得できませんでした。"})
        return
    }
    // 認証コードでトークンを取得
    token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Google認証に失敗しました。"})
        return
    }
    // Google APIからユーザー情報を取得
    resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "ユーザー情報の取得に失敗しました。"})
        return
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "ユーザー情報の読み取りに失敗しました。"})
        return
    }
    var userInfo struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    if err := json.Unmarshal(data, &userInfo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "ユーザー情報の解析に失敗しました。"})
        return
    }
    // 既存ユーザーがいれば取得、いなければ新規作成
    var user models.User
    if err := config.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
        // 新規ユーザー作成（パスワードは空）
        user = models.User{Name: userInfo.Name, Email: userInfo.Email, Password: ""}
        config.DB.Create(&user)
    }
    // JWTトークン発行
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := authToken.SignedString([]byte(config.JWTSecret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "トークン生成に失敗しました。"})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "message": "Googleログインに成功しました。",
        "token":   tokenString,
    })
}
// Logout はユーザーのログアウトを処理します。
func Logout(c *gin.Context) {
	// ログアウトはクライアント側でトークンを破棄するだけで十分です。
	// サーバー側では特に処理は必要ありません。
	c.JSON(http.StatusOK, gin.H{"message": "ログアウトしました。"})
}	
