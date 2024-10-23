package model

import "gorm.io/gorm"

type TssRequest struct {
	gorm.Model

	RequestID     string `gorm:"size:255;not null;default:'';uniqueIndex:ux_request_id"`
	Type          int64  `gorm:"not null;default:0"`
	Status        int64  `gorm:"not null;default:0"`
	TargetGroupID string `gorm:"size:255;not null;default:'';index:,length:8"`
	SourceGroupID string `gorm:"size:255;not null;default:'';index:,length:8"`
	UserID        string `gorm:"size:255;not null;default:'';index:,length:8"`
	VaultID       string `gorm:"size:255;not null;default:''"`
}
