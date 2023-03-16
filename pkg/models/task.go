package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	Text      string         `json:"text" binding:"required"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
