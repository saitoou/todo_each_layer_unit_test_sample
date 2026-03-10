package entity

import "time"

type Todo struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	// UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
