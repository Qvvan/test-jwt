package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

const (
	pgErrDuplicateCode = "SQLSTATE 23505"
)

var ErrDuplicate = errors.New("duplicate record")

type BaseRepo struct {
	db    *sqlx.DB
	table string
}

func (r *BaseRepo) getQuery(id string) (string, error) {
	qb := goqu.From(r.table)
	qb = qb.Where(goqu.I("id").Eq(id))

	query, _, err := qb.ToSQL()
	if err != nil {
		return "", fmt.Errorf("can't build query to get: %w", err)
	}

	return query, nil
}

func (r *BaseRepo) GetID(id string, getModel interface{}) error {
	q, err := r.getQuery(id)
	if err != nil {
		return err
	}
	if err := r.db.Get(getModel, q); err != nil {
		return fmt.Errorf("can't get model: %w", err)
	}

	return nil
}

func (r *BaseRepo) update(updateModel interface{}, id string) error {
	qu, _, err := goqu.Update(r.table).Set(
		updateModel,
	).Returning("id").Where(goqu.I("id").Eq(id)).ToSQL()
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(qu); err != nil {
		if strings.Contains(err.Error(), pgErrDuplicateCode) {
			return ErrDuplicate
		}
		return fmt.Errorf("error can't update updateModel %w", err)
	}

	return nil
}
