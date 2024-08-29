package models

import "time"

type Note struct {
	Id        int       `json:"id,omitempty" db:"id"`
	UserId    int       `json:"user_id,omitempty" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Text      string    `json:"text" db:"text"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}