package models

import "github.com/jinzhu/gorm"

// Role represents the role model
type Role struct {
	gorm.Model
	Name        string       `gorm:"size:255;not null;unique" json:"name"`
	Description string       `gorm:"size:255;" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}
