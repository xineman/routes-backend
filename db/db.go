package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DB_URL = "postgres://postgres:postgres@localhost/postgres"

var DbPool *pgxpool.Pool

func Init() {
	var err error
	DbPool, err = pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
}
