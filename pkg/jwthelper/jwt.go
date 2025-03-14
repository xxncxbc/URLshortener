package jwthelper

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	Secret string
}

type JWTData struct {
	Email  string
	UserId uint
}

func NewJWT(secret string) *JWT {
	return &JWT{secret}
}

func (j *JWT) Create(data JWTData, exp time.Time) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   data.Email,
		"exp":     exp.Unix(),
		"user_id": data.UserId,
	})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"].(string)
	userId := uint(t.Claims.(jwt.MapClaims)["user_id"].(float64))
	return t.Valid, &JWTData{email, userId}
}
