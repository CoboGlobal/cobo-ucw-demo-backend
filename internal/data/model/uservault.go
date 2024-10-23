package model

import "gorm.io/gorm"

type UserVault struct {
	gorm.Model
	UserID  string `gorm:"size:255;not null;default:'';uniqueIndex:ux_user_id_vault_id"`
	VaultID string `gorm:"size:255;not null;default:'';uniqueIndex:ux_user_id_vault_id"`
}
