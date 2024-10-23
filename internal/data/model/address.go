package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	WalletID string `gorm:"size:255;not null;default:'';index:ix_wallet_id_chain_id"`
	ChainID  string `gorm:"size:255;not null;default:'';index:ix_wallet_id_chain_id"`
	Address  string `gorm:"size:255;not null;default:''"`
	Path     string `gorm:"size:255;not null;default:''"`
	PubKey   string `gorm:"size:255;not null;default:''"`
	Encoding string `gorm:"size:255;not null;default:''"`
}
