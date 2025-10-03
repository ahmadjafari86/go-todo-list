package models

import "time"

type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Email        string    `gorm:"type:text;not null;unique" json:"email" binding:"required,email"`
    PasswordHash string    `gorm:"type:text;not null" json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}
