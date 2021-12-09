package model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	Chat string
}

type ChatRequest struct {
	Chat string
}


type ChatResponse struct {
	ID        uint `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time  `json:"deleted_at"`
	Chat string `json:"chat"`
	Status int `json:"status"`
}
