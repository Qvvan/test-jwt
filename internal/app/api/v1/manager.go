package v1

import (
	"log/slog"

	"github.com/qvvan/test-jwt/internal/app/repository"
)

type Manager struct {
	factory *repository.Factory
	log     *slog.Logger
}

func NewManager(
	factory *repository.Factory,
	log *slog.Logger,
) *Manager {
	return &Manager{
		factory: factory,
		log:     log,
	}
}
