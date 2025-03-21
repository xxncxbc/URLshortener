package user

import (
	"URLshortener/internal/link"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"index"`
	Links    []link.Link
}
