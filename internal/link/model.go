package link

import (
	"URLshortener/internal/stat"
	"gorm.io/gorm"
	"math/rand"
)

type Link struct {
	gorm.Model
	Url    string      `json:"url"`
	Hash   string      `json:"hash" gorm:"uniqueIndex"`
	Stats  []stat.Stat `json:"-" gorm:"constraints:OnUpdate:CASCADE,OnDelete:Set NULL;"`
	UserId uint        `json:"user_id" gorm:"foreignKey:UserId;references:ID"`
}

func NewLink(url string, userId uint) *Link {
	link := &Link{
		Url:    url,
		UserId: userId,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = RandStringRunes(10)
}

var letterRunes = []rune("abcdefghijklmnoprstuvwxyzABCDEFGHIJKLMNOPRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
