package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"unique;size:50"`
	Password string `json:"password" binding:"required" gorm:"size:100"`
}
