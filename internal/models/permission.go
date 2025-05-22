package models

import "github.com/jinzhu/gorm"

// Permission represents the permission model
type Permission struct {
	gorm.Model
	Name        string `gorm:"size:255;not null;unique" json:"name"`
	Description string `gorm:"size:255;" json:"description"`
	Resource    string `gorm:"size:255;not null;" json:"resource"`
	Action      string `gorm:"size:255;not null;" json:"action"`
}
