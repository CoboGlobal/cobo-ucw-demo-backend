package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model

	VaultID   string `gorm:"size:255;not null;default:'';index,length:8"`
	GroupID   string `gorm:"size:255;not null;default:'';uniqueIndex:ux_group_id"`
	GroupType int64  `gorm:"not null;default:0"`
}
