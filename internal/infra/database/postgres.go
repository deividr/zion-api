package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return dbpool
}
