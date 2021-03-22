package models

import (
	"time"

	"gorm.io/gorm"
)

// Post holds all the metadata and content of a post
type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// PostData the user-editable part of the Post
	PostData `json:",flow"`
}

// PostData includes all the user-editable fields
type PostData struct {
	Title  *string `json:"title"`
	Author *string `json:"author"`
	Text   *string `json:"text"`
}
