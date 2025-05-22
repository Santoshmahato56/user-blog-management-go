package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user model
type User struct {
	gorm.Model
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Email     string `gorm:"size:255;not null;unique" json:"email"`
	Password  string `gorm:"size:255;not null;" json:"password,omitempty"`
	FirstName string `gorm:"size:255;" json:"first_name,omitempty"`
	LastName  string `gorm:"size:255;" json:"last_name,omitempty"`
	RoleID    uint   `gorm:"not null;" json:"role_id"`
	Role      Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// BeforeSave is a hook that runs before saving the user
func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ValidatePassword checks if the provided password is correct
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
