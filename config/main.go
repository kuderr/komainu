package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(dataSourceName string) {
	// var err error
	// DB, err = sql.Open("postgres", dataSourceName)
	// if err != nil {
	// 	log.Panic(err)
	// }

	var err error
	DB, err = pgxpool.Connect(context.Background(), dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(context.Background()); err != nil {
		log.Panic(err)
	}
}

func CloseDB() {
	DB.Close()
}
