package models

import "time"

type Note struct {
	Id        int       `json:"id,omitempty" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	UserName  string    `json:"username" db:"username"`
	Title     string    `json:"title" db:"title"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}