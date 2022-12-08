package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Hash      string    `json:"hash"`
	Salt      string    `json:"salt"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "admin_user"
}
