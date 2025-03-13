package middleware

import (
	"URLshortener/configs"
	"URLshortener/pkg/jwthelper"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeAuthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	return
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(AuthHeader, "Bearer ") {
			writeAuthed(w)
			return
		}
		token := strings.TrimPrefix(AuthHeader, "Bearer ")
		isValid, data := jwthelper.NewJWT(config.Auth.AccessSecret).Parse(token)
		if !isValid {
			writeAuthed(w)
			return
		}
		//тут приходится создавать новый контекст и новый запрос и передавать
		//значение авторизации в него, чтобы обработчики ниже уже знали о статусе
		//авторизации пользователя
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		//передаем новый запрос в обработку
		next.ServeHTTP(w, req)
	})
}
