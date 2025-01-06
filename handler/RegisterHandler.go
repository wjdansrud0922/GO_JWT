package handler

import (
	"GO_JWT/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		var User entity.User
		var DbUser entity.User
		if err := c.ShouldBindJSON(&User); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":       "회원가입 실패 | 입력 정보가 잘못되었습니다.",
				"error msg": err.Error(),
			})
			return
		}

		if err := db.Where("username = ?", User.Username).First(&DbUser).Error; err == nil { //존재하면 nil 반환(First)
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "회원가입 실패 | 유저 이름 중복",
			})
			return
		}

		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "회원가입 실패, 암호화 오류",
				"err": err.Error(),
			})
			return
		}

		User.Password = string(bcryptPassword)
		db.Create(&User)
		c.JSON(http.StatusOK, gin.H{
			"msg": "회원가입 성공",
		})
		return
	}
}
