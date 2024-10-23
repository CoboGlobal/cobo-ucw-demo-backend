package model

import (
	"gorm.io/gorm"
)

type GroupNode struct {
	gorm.Model
	NodeID     string `gorm:"size:255;not null;default:'';uniqueIndex:ux_group_id_node_id,priority:2"`
	GroupID    string `gorm:"size:255;not null;default:'';uniqueIndex:ux_group_id_node_id,priority:1"`
	HolderName string `gorm:"size:255;not null;default:''"`
	UserID     string `gorm:"size:255;not null;default:'';index:,length:8"`
}
