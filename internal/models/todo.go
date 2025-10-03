package models

import "time"

type Todo struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"type:text;not null" json:"title" binding:"required"`
    Completed bool      `gorm:"not null" json:"completed"`
    OwnerID   uint      `gorm:"not null" json:"owner_id"`
    CreatedAt time.Time `json:"created_at"`
}
