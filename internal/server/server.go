package server

import (
	"log"
	"user-management/internal/auth"
	"user-management/internal/chat"
	"user-management/internal/config"
	db "user-management/internal/db/generated"
	"user-management/internal/user"
	"user-management/pkg/redisx"
	"user-management/pkg/validation"
)


func Start(pool *db.Queries) error {
	SetupLogger()
	RegisterMetrics()

	if err := validation.InitValidator(); err != nil {
		log.Fatal("Failed to initialize validator: ", err)
	}

	rdb := redisx.NewClient(config.Load.REDIS_URL, "", 0)
	cache := redisx.NewMessageCache(rdb)

	authService := auth.NewService(pool)
	userService := user.NewService(pool)
	chatService := chat.NewService(pool, cache)

	router := SetupRouter(authService, userService, chatService, rdb)

	port := config.Load.PORT
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port: %s", port)
	return router.Run(":" + port)
}