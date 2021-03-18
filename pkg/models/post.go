package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// PostData the user-editable part of the Post
	PostData `json:",flow"`
}

type PostData struct {
	Title  *string `json:"title"`
	Author *string `json:"author"`
	Text   *string `json:"text"`
}
