package model

import "gorm.io/gorm"

type UserNode struct {
	gorm.Model

	UserID string `gorm:"size:255;not null;default:'';uniqueIndex:ux_user_id_node_id"`
	NodeID string `gorm:"size:255;not null;default:'';uniqueIndex:ux_user_id_node_id"`
	Status int64  `gorm:"not null;default:0"`
}
