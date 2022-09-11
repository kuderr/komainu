package main

import (
	autherApi "auther/api/auther"
	"auther/config"
	"auther/internal/auther"
	"auther/internal/database"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Read configuration
	cfg, err := config.Read("env", ".env", ".")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiates the database
	postgres, err := database.NewPostgres(cfg.PostgresUrl)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer postgres.DB.Close()

	queries := database.New(postgres.DB)
	dbAuthStorage := auther.NewDatabaseAuthInfoStorage(queries)
	builder := auther.NewBuilder(dbAuthStorage)
	accesses, clients, err := builder.BuildAccessMap()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiates the author service
	authInfo := auther.NewAuthInfo(accesses, clients)
	authService := autherApi.NewService(authInfo, cfg.Secret)

	// Register our service handlers to the router
	router := gin.Default()
	authService.RegisterHandlers(router)

	go syncAuthInfo(builder, authInfo)

	// Start the server
	router.Run()
}

func syncAuthInfo(builder *auther.Builder, authInfo *auther.AuthInfo) {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("sync auth info")
			// TODO: Maybe pass pointers
			accesses, clients, err := builder.BuildAccessMap()
			if err != nil {
				log.Println("Error: ", err.Error())
			}
			authInfo.Update(accesses, clients)
		}
	}
}
