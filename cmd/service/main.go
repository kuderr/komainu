package main

import (
	"auther/api/auther"
	"auther/config"
	"auther/internal/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Read configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(cfg)

	// Instantiates the database
	postgres, err := database.NewPostgres(cfg.PostgresUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiates the author service
	queries := database.New(postgres.DB)
	authService := auther.NewService(queries, cfg.Secret)

	// Register our service handlers to the router
	router := gin.Default()
	authService.RegisterHandlers(router)

	// Start the server
	router.Run()
}
