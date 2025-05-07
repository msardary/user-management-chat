package db

import (
	"context"
	"log"

	"user-management/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect() (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}

	return pool, nil

}

func SetPool(p *pgxpool.Pool) {
	pool = p
}

func Ping(ctx context.Context) error {
	if pool == nil {
		return nil
	}
	return pool.Ping(ctx)
}