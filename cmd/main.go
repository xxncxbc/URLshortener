package main

import (
	"URLshortener/configs"
	"URLshortener/internal/auth"
	"URLshortener/internal/link"
	"URLshortener/internal/user"
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
	userRepository := user.NewUserRepository(database)
	//services
	authService := auth.NewAuthService(userRepository)
	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
	})

	//middlewares снизу вверх
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	//сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}
	fmt.Println("Server is listening on port 8080")
	// запуск сервера
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
