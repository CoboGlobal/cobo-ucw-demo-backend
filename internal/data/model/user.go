package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	UserID string `gorm:"size:255;not null;default:'';uniqueIndex:ux_user_id"`
	Email  string `gorm:"size:255;not null;default:'';uniqueIndex:ux_email"`
}
