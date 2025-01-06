package main

import (
	"GO_JWT/db"
	"GO_JWT/entity"
	"GO_JWT/handler"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func main() {
	db := db.InitDb()

	jwt, err := jwt.New(InitParams(db))
	if err != nil {
		panic("jwt 오류")
	}
	router := gin.Default()

	router.POST("/register", handler.RegisterHandler(db))
	router.POST("/login", handler.LoginHandler(db, jwt))

	authRouter := router.Group("/auth")
	authRouter.Use(jwt.MiddlewareFunc())
	{
		authRouter.GET("/user", handler.TestHandler)
	}

	router.Run(":8080")
}

var identityKey string = "id"

func InitParams(db *gorm.DB) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:           "test",
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour * 3,
		IdentityKey:     identityKey,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(db),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		claims := jwt.ExtractClaims(c)
		if v, ok := data.(*entity.User); ok && v.Username == claims[identityKey].(string) {
			return true
		}
		return false
	}
}

func authenticator(db *gorm.DB) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var User entity.User
		var DbUser entity.User
		if err := c.ShouldBindJSON(&User); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		if err := db.Where("username = ?", User.Username).First(&DbUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, jwt.ErrFailedAuthentication
			}
			return nil, err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(User.Password)); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
		return &entity.User{
			Username: DbUser.Username,
		}, nil

	}
}

func identityHandler() func(*gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &entity.User{
			Username: claims[identityKey].(string),
		}
	}

}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*entity.User); ok {
			return jwt.MapClaims{
				identityKey: v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}
