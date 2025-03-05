package auth

import (
	"URLshortener/configs"
	"URLshortener/pkg/req"
	"URLshortener/pkg/res"
	"fmt"
	"log"
	"net/http"
)

// зависимости конфига
type AuthHandlerDeps struct {
	*configs.Config
}
type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Login")
		var payload *LoginRequest
		payload, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println(*payload)
		data := LoginResponse{
			Token: "lol",
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Register")
		var payload *RegisterRequest
		payload, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println(*payload)
		res.Json(w, *payload, http.StatusOK)
	}
}
