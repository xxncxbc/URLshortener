package main

import (
	"URLshortener/configs"
	"URLshortener/internal/auth"
	"URLshortener/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	server := &http.Server{Addr: ":8080", Handler: router}
	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
