package auth

import (
	"URLshortener/configs"
	"URLshortener/pkg/jwthelper"
	"URLshortener/pkg/req"
	"URLshortener/pkg/res"
	"net/http"
	"time"
)

const (
	AccessTokenLifeSpan  = time.Hour * 3
	RefreshTokenLifeSpan = time.Hour * 24 * 7
)

// зависимости конфига
type AuthHandlerDeps struct {
	*AuthService
	*configs.Config
}
type AuthHandler struct {
	AuthService *AuthService
	Config      *configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthService: deps.AuthService,
		Config:      deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/refresh", handler.Refresh())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		email, userId, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		accessToken, err := jwthelper.NewJWT(handler.Config.Auth.AccessSecret).Create(jwthelper.JWTData{
			Email:  email,
			UserId: userId,
		},
			time.Now().Add(AccessTokenLifeSpan))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		refreshToken, err := jwthelper.NewJWT(handler.Config.Auth.RefreshSecret).Create(jwthelper.JWTData{
			Email:  email,
			UserId: userId,
		}, time.Now().Add(RefreshTokenLifeSpan))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		res.Json(w, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		email, userId, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		accessToken, err := jwthelper.NewJWT(handler.Config.Auth.AccessSecret).Create(jwthelper.JWTData{
			Email:  email,
			UserId: userId,
		},
			time.Now().Add(AccessTokenLifeSpan))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		refreshToken, err := jwthelper.NewJWT(handler.Config.Auth.RefreshSecret).Create(jwthelper.JWTData{
			Email:  email,
			UserId: userId,
		},
			time.Now().Add(RefreshTokenLifeSpan))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := RegisterResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		res.Json(w, data, http.StatusCreated)
	}
}

func (handler *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RefreshRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		isValid, parsedToken := jwthelper.NewJWT(handler.Config.Auth.RefreshSecret).Parse(body.RefreshToken)
		if !isValid {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}
		accessToken, err := jwthelper.NewJWT(handler.Config.Auth.AccessSecret).Create(jwthelper.JWTData{
			Email:  parsedToken.Email,
			UserId: parsedToken.UserId,
		},
			time.Now().Add(AccessTokenLifeSpan))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		refreshToken, err := jwthelper.NewJWT(handler.Config.Auth.AccessSecret).Create(jwthelper.JWTData{
			Email:  parsedToken.Email,
			UserId: parsedToken.UserId,
		},
			time.Now().Add(RefreshTokenLifeSpan))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := RefreshResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		res.Json(w, data, http.StatusOK)
	}
}
