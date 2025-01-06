package handler

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	username := claims["id"].(string)

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "testHandler jwt 존재 x"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "good"})
}
