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
	*AuthService
	*configs.Config
}
type AuthHandler struct {
	*AuthService
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthService: deps.AuthService,
		Config:      deps.Config,
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
		email, err := handler.AuthService.Login(payload.Email, payload.Password)
		fmt.Println(email)
		data := LoginResponse{
			Token: "lol",
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Register")
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			log.Println(err.Error())
			return
		}
		handler.AuthService.Register(body.Email, body.Password, body.Name)
		res.Json(w, *body, http.StatusOK)
	}
}
