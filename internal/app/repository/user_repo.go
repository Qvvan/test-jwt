package repository

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/qvvan/test-jwt/internal/app/models"
)

const (
	UsersTName = "users"
)

type UserRepo struct {
	BaseRepo
	log *slog.Logger
}

func NewUserRepo(db *sqlx.DB, log *slog.Logger) *UserRepo {
	repo := &UserRepo{}
	repo.db = db
	repo.table = UsersTName
	repo.log = log
	return repo
}

func (r *UserRepo) GetID(id string) (*models.User, error) {
	getModel := new(models.User)
	q, err := r.getQuery(id)
	if err != nil {
		return nil, err
	}
	if err := r.db.Get(getModel, q); err != nil {
		return nil, fmt.Errorf("can't get user model: %w", err)
	}

	return getModel, nil
}

func (r *UserRepo) Update(user *models.User) error {
	if err := r.BaseRepo.update(user, user.ID); err != nil {
		r.log.Error("failed to update user", slog.Any("err", err))
		return err
	}

	return nil
}
