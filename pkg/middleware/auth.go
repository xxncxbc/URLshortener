package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(AuthHeader, "Bearer ")
		fmt.Println(token)
		next.ServeHTTP(w, r)
	})
}
