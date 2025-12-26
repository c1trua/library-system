package models

import (
	"time"
)

type BorrowRecord struct {
	ID         int        `gorm:"primaryKey" json:"id" example:"1"`
	UserID     int        `json:"user_id" example:"1"`
	BookID     int        `json:"book_id" example:"1"`
	BorrowedAt time.Time  `json:"borrowed_at" example:"2024-01-15T10:30:00Z"`
	DueDate    time.Time  `json:"due_date" example:"2024-02-15T10:30:00Z"`
	ReturnedAt *time.Time `json:"returned_at,omitempty" example:"2024-01-18T09:15:00Z"`
}
