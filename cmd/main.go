package main

import (
	"URLshortener/configs"
	"URLshortener/internal/auth"
	"URLshortener/internal/link"
	"URLshortener/internal/stat"
	"URLshortener/internal/user"
	"URLshortener/pkg/db"
	"URLshortener/pkg/event"
	"URLshortener/pkg/middleware"
	"fmt"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//Repositories
	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)
	//services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})
	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})
	//middlewares снизу вверх
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	go statService.AddClick()
	return stack(router)
}

func main() {
	app := App()
	//сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}
	fmt.Println("Server is listening on port 8080")
	// запуск сервера
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
