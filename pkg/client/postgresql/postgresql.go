package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const defaultTimeoutSec = 30

func NewClient(ctx context.Context, pgDSN string, maxRetries int, retryDelay time.Duration) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeoutSec*time.Second)
	defer cancel()

	for i := 0; i < maxRetries; i++ {
		var err error
		db, err := sqlx.ConnectContext(ctx, "pgx", pgDSN)
		if err == nil {
			log.Println("Successfully connected to database.")
			err = db.PingContext(ctx)
			if err == nil {
				return db, nil
			}
			log.Println("Ping failed: ", err)
			db.Close()
		}
		log.Printf("Attempt %d failed: %v. Retrying in %v...\n", i+1, err, retryDelay)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelay):
		}
	}
	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}
