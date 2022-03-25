package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBInterface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type DBManager struct {
	Pool DBInterface
}

func Connect() *DBManager {
	connString := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println("Connection to postgres failed")
		return nil
	}
	fmt.Println("Successful connection to postgres")
	return &DBManager{Pool: pool}
}

func Disconnect(db *DBManager) {
	db.Pool.Close()
}
