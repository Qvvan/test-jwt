package repository

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/qvvan/test-jwt/internal/app/models"
	errors "github.com/qvvan/test-jwt/pkg/client/postgresql/utils"
)

const (
	UsersTName = "users"
)

type UserRepo struct {
	BaseRepo
	log *slog.Logger
}

func NewUserRepo(db *pgxpool.Pool, log *slog.Logger) *UserRepo {
	repo := &UserRepo{}
	repo.db = db
	repo.table = UsersTName
	repo.log = log
	return repo
}

func (r *UserRepo) GetID(ctx *gin.Context, id string) (*models.User, *errors.CustomError) {
	newModel := new(models.User)
	if err := r.BaseRepo.GetID(ctx, id, newModel); err != nil {
		r.log.Error("failed to get user", slog.Any("err", err))
		return nil, err
	}

	return newModel, nil
}

func (r *UserRepo) Create(ctx *gin.Context, user *models.User) (*models.User, *errors.CustomError) {
	id, err := r.BaseRepo.create(ctx, user)
	if err != nil {
		r.log.Error("failed to create user", slog.Any("err", err))
		return nil, err
	}
	user.UserID = id
	return user, nil
}

func (r *UserRepo) Update(ctx *gin.Context, user *models.User) error {
	if err := r.BaseRepo.update(ctx, user, user.UserID); err != nil {
		r.log.Error("failed to update user", slog.Any("err", err))
		return err
	}

	return nil
}

func (r *UserRepo) Delete(user *models.User) (string, error) {
	return "User deleted successfully", nil
}
