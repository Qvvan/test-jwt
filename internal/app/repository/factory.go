package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Factory struct {
	UserRepo *UserRepo
}

func NewFactory(
	db *pgxpool.Pool,
	log *slog.Logger,
) *Factory {
	return &Factory{
		UserRepo: NewUserRepo(db, log),
	}
}
