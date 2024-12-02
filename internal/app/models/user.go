package models

import "time"

type User struct {
	UserID       string    `db:"user_id" goqu:"defaultifempty,skipinsert,skipupdate"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at,omitempty"`
}
