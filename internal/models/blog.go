package models

import "github.com/jinzhu/gorm"

// Blog represents the blog post model
type Blog struct {
	gorm.Model
	Title     string `gorm:"size:255;not null;" json:"title"`
	Content   string `gorm:"type:text;not null;" json:"content"`
	Published bool   `gorm:"default:false" json:"published"`
	UserID    uint   `gorm:"not null;" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
