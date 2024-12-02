package repository

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	errors "github.com/qvvan/test-jwt/pkg/client/postgresql/utils"
)

type BaseRepo struct {
	db    *pgxpool.Pool
	table string
}

func (r *BaseRepo) create(ctx *gin.Context, newModel interface{}) (string, *errors.CustomError) {
	query, _, err := goqu.Insert(r.table).Rows(newModel).Returning("id").ToSQL()
	if err != nil {
		return "", errors.NewQueryError(r.table, "query build", err)
	}

	var id string
	if err := r.db.QueryRow(ctx, query).Scan(&id); err != nil {
		parsedErr := errors.ParsePostgresError(err)
		return "", errors.NewCreateError(r.table, "scan result", parsedErr)
	}

	return id, nil
}

func (r *BaseRepo) getQuery(id string) (string, *errors.CustomError) {
	qb := goqu.From(r.table).Where(goqu.I("id").Eq(id))
	query, _, err := qb.ToSQL()
	if err != nil {
		return "", errors.NewQueryError(r.table, "query build for get", err)
	}
	return query, nil
}

func (r *BaseRepo) GetID(ctx *gin.Context, id string, newModel interface{}) *errors.CustomError {
	query, err := r.getQuery(id)
	if err != nil {
		return err
	}

	if err := r.db.QueryRow(ctx, query).Scan(&newModel); err != nil {
		parsedErr := errors.ParsePostgresError(err)
		return errors.NewCreateError(r.table, "scan result", parsedErr)
	}

	return nil
}

func (r *BaseRepo) update(ctx *gin.Context, updateModel interface{}, id string) error {
	qu, _, err := goqu.Update(r.table).Set(
		updateModel,
	).Returning("id").Where(goqu.I("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(ctx, qu); err != nil {
		if strings.Contains(err.Error(), errors.PGErrDuplicateCode) {
			return errors.ErrDuplicate
		}
		return fmt.Errorf("error can't update updateModel %w", err)
	}

	return nil
}
