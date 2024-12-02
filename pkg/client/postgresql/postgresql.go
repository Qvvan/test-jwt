package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, dsn string) (pool *pgxpool.Pool, err error) {
	err = DoWithAttemps(func() error {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pxCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		pool, err = pgxpool.ConnectConfig(attemptCtx, pxCfg)
		if err != nil {
			log.Println("failed to connect to database: ", err.Error())
			return err
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		return nil, fmt.Errorf("all attempts exceeded: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return pool, nil
}

func DoWithAttemps(fn func() error, maxAttempts int, maxDelay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err == nil {
			return nil
		}

		maxAttempts--
		time.Sleep(maxDelay)
	}

	return err
}
