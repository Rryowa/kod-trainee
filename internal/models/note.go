package models

import "time"

type Note struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}