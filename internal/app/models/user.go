package models

import "time"

type User struct {
	UserID    string    `db:"user_id" goqu:"defaultifempty,skipinsert,skipupdate"`
	CreatedAt time.Time `db:"created_at,omitempty"`
}
