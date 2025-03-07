package main

import (
	"URLshortener/configs"
	"URLshortener/internal/auth"
	"URLshortener/internal/link"
	"URLshortener/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	//Repositories
	linkRepository := link.NewLinkRepository(database)
	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	server := &http.Server{Addr: ":8080", Handler: router}
	fmt.Println("Server is listening on port 8080")
	server.ListenAndServe()
}
