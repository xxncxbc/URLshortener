package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			//каждый раз оборачиваем next в новый middleware
			next = middlewares[i](next)
		}
		return next
	}
}
