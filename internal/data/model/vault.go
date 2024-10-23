package model

import "gorm.io/gorm"

type Vault struct {
	gorm.Model

	VaultID     string `gorm:"size:255;not null;default:'';uniqueIndex:ux_vault_id"`
	Name        string `gorm:"size:255;not null;default:''"`
	MainGroupID string `gorm:"size:255;not null;default:''"`
	ProjectID   string `gorm:"size:255;not null;default:''"`
	Status      int64  `gorm:"not null;default:0"`
}
