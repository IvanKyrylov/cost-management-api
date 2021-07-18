package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewClient(ctx context.Context, host, port, user, password, dbname, sslmode string) (client *sql.DB, err error) {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err = sql.Open("postgres", connect)
	if err != nil {
		return nil, fmt.Errorf("failed to create client to postgres due to error %w", err)
	}
	err = client.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client to postgres due to error %w", err)
	}
	return client, nil
}
