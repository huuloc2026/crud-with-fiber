package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // "-" prevents password from being sent in JSON responses
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
