package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model

	VaultID  string `gorm:"size:255;not null;default:''"`
	UserID   string `gorm:"size:255;not null;default:''"`
	WalletID string `gorm:"size:255;not null;default:'';uniqueIndex:ux_wallet_id"`
	Name     string `gorm:"size:255;not null;default:''"`
}
