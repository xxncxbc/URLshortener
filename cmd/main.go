package main

import (
	"URLshortener/configs"
	"URLshortener/internal/auth"
	"URLshortener/internal/link"
	"URLshortener/pkg/db"
	"URLshortener/pkg/middleware"
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

	//middlewares снизу вверх
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	server := &http.Server{Addr: ":8080",
		Handler: stack(router),
	}
	fmt.Println("Server is listening on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
