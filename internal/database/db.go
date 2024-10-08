package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Seven11Eleven/auth_service_medods/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgreSQLConnection(env *config.Env) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("unable to connect to database : %v", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("unable to ping database : %v", err)
	}

	log.Println("connected to database succdessfully")

	return conn
}

func ClosePostgreSQLConnection(conn *pgx.Conn) {
	if conn == nil {
		return
	}

	err := conn.Close(context.Background())
	if err != nil {
		log.Fatalf("error closing database conn: %v", err)

	}
	log.Println("connection to database is closed!")
}