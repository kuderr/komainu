package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"auther/auther"
	"auther/config"

	"github.com/julienschmidt/httprouter"
)

func main() {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	config.InitDB(dbUrl)

	router := httprouter.New()

	router.POST("/auth/", auther.CheckAccess)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server started at http://127.0.0.1:5000")

	log.Fatal(srv.ListenAndServe())
}

func auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
