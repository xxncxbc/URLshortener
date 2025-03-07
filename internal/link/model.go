package link

import (
	"gorm.io/gorm"
	"math/rand"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
}

func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: RandStringRunes(10),
	}
}

var letterRunes = []rune("abcdefghijklmnoprstuvwxyzABCDEFGHIJKLMNOPRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
