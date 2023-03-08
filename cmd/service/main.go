package main

import (
	api "checker/api/checker"
	"checker/config"
	checker "checker/internal/core"
	"checker/internal/database"
	"log"
	"net/http"
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
	dbAuthStorage := checker.NewDatabaseAuthInfoStorage(queries)
	builder := checker.NewBuilder(dbAuthStorage)
	accesses, clients, err := builder.BuildAccessMap()
	if err != nil {
		log.Fatal(err.Error())
	}

	// For debugging access map
	// accessMapJSON, _ := json.MarshalIndent(accesses, "", "    ")
	// log.Println(string(accessMapJSON))
	// clientsJSON, _ := json.MarshalIndent(clients, "", "    ")
	// log.Println(string(clientsJSON))

	// Instantiates the author service
	authInfo := checker.NewAuthInfo(accesses, clients)
	authService := api.NewService(authInfo, cfg.JWTPublicKey)

	// Register our service handlers to the router
	router := gin.Default()
	router.GET("/livez", healthCheck)
	router.GET("/readyz", healthCheck)

	authService.RegisterHandlers(router)

	go syncAuthInfo(builder, authInfo)

	// Start the server
	router.Run(":5000")
}

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func syncAuthInfo(builder *checker.Builder, authInfo *checker.AuthInfo) {
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
