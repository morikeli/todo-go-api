package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dbURL string) (*pgxpool.Pool, error) {
	var ctx context.Context = context.Background()

	var config *pgxpool.Config
	var err error

	config, err = pgxpool.ParseConfig(dbURL)

	if err != nil {
		log.Printf("[ERROR]: Unable to parse database url: %v\n", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("[ERROR]: Unable to create a connection pool: %v\n", err)
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Printf("Unable to ping db: %v", err)
		pool.Close()

		return nil, err
	}

	log.Println("[SUCCESS]: Successfully connected to PostgreSQL db!")

	return pool, nil
}
