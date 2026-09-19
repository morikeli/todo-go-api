package models

import "time"

type Task struct {
	Id int `json:"id" db:"id"`
	UserId string `json:"user_id" db:"user_id"`
	Title string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	Completed bool `json:"is_completed" db:"completed"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

}