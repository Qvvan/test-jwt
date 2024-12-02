package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Factory struct {
	UserRepo *UserRepo
}

func NewFactory(
	db *sqlx.DB,
	log *slog.Logger,
) *Factory {
	return &Factory{
		UserRepo: NewUserRepo(db, log),
	}
}
