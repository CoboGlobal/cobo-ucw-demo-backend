package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	TransactionID string          `gorm:"size:255;not null;default:'';uniqueIndex:ux_transaction_id"`
	WalletID      string          `gorm:"size:255;not null;default:''"`
	Type          int64           `gorm:"not null;default:0"`
	Chain         string          `gorm:"size:255;not null;default:''"`
	Amount        decimal.Decimal `gorm:"size:255;not null;default:''"`
	From          string          `gorm:"size:255;not null;default:''"`
	To            string          `gorm:"size:255;not null;default:''"`
	TxHash        string          `gorm:"size:255;not null;default:''"`
	Fee           Fee             `gorm:"type:json;serializer:json"`
	Status        int64           `gorm:"not null;default:0"`
	ExternalID    string          `gorm:"size:255;uniqueIndex:ux_transactions_external_id"` // cobo transaction id
	TokenID       string          `gorm:"size:255;not null;default:''"`
	BlockNum      int64           `gorm:"not null;default:0"`
	Extra         Extra           `gorm:"type:json;serializer:json"`
	SubStatus     int64           `gorm:"not null;default:0"`
	UserID        string          `gorm:"size:255;not null;default:'';index:,length:8"`
}

type Extra struct {
	FailedReason        string `json:"failed_reason,omitempty"`
	ConfirmedNum        int64  `json:"confirmed_num,omitempty"`
	ConfirmingThreshold int64  `json:"confirming_threshold,omitempty"`
	BlockHash           string `json:"block_hash,omitempty"`
	Description         string `json:"description,omitempty"`
	Nonce               int64  `json:"nonce,omitempty"`
	RawTx               string `json:"raw_tx,omitempty"`
}

type Fee struct {
	TokenID        string          `json:"token_id"`
	GasPrice       decimal.Decimal `json:"gas_price"`
	GasLimit       decimal.Decimal `json:"gas_limit"`
	FeePerByte     decimal.Decimal `json:"fee_per_byte"`
	FeeAmount      decimal.Decimal `json:"fee_amount"`
	Level          int64           `json:"level"`
	MaxPriorityFee decimal.Decimal `json:"max_priority_fee"`
	MaxFee         decimal.Decimal `json:"max_fee"`
	FeeUsed        decimal.Decimal `json:"fee_used"`
}
