package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Hidden in JSON response
	CreatedAt time.Time `json:"created_at"`
}