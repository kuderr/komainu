package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Panic(err)
	}

	// Try connecting to the database
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(60 * time.Second)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("connection timeout")

		case <-ticker.C:
			if err := db.Ping(context.Background()); err == nil {
				return &Postgres{DB: db}, nil
			}
		}
	}
}
