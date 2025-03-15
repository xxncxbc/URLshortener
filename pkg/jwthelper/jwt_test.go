package jwthelper

import (
	"testing"
	"time"
)

func TestJWT_CreateAndParse(t *testing.T) {
	const (
		email       = "a@a.com"
		userId uint = 1
	)
	jwtService := NewJWT("f4c613b58c39b8d35b9b827d51bc9a7149ea58b1b1d8b36c6603daff0770d899")
	token, err := jwtService.Create(JWTData{
		Email:  email,
		UserId: userId,
	},
		time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err.Error())
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal(err.Error())
	}
	if data.Email != email || data.UserId != userId {
		t.Fatal(data.Email, data.UserId)
	}
}
