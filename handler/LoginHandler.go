package handler

import (
	"GO_JWT/entity"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func LoginHandler(db *gorm.DB, jwtMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestUser entity.User
		var dbUser entity.User

		if err := c.ShouldBindJSON(&requestUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "입력값 오류"})
			log.Println(err.Error())
			return
		}

		if err := db.Where("username = ?", requestUser.Username).First(&dbUser).Error; err != nil { // 유저 이름에 맞는 패스워드가 존재하지 않으면 오류 즉 유저이름이 저장되어있지 않음
			c.JSON(http.StatusNotFound, gin.H{"msg": "아이디 또는 패스워드가 틀렸습니다."})
			log.Println(err.Error())
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(requestUser.Password)); err != nil { // 해당 유저의 패스워드와 유저이름으로 조회한 암호화된 패스워드가 다르면 오류 즉 비밀번호가 틀림
			c.JSON(http.StatusNotFound, gin.H{"msg": "아이디 또는 패스워드가 틀렸습니다."})
			log.Println(err.Error())
			return
		}

		//토큰

		token, _, err := jwtMiddleware.TokenGenerator(&requestUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "토큰 생성 실패"})
			log.Println(err.Error())
			return
		}

		// 토큰 반환
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
