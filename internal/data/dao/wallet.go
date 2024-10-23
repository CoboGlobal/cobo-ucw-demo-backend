package dao

import (
	"context"

	"cobo-ucw-backend/internal/biz/wallet"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm/clause"
)

func NewWallet(db *database.Data) *Wallet {
	return &Wallet{db}
}

type Wallet struct {
	*database.Data
}

func (d *Wallet) Update(ctx context.Context, data wallet.Wallet) error {
	return d.DB.WithContext(ctx).Model(&data.Wallet).Updates(data.Wallet).Error
}

func (d *Wallet) Save(ctx context.Context, data *wallet.Wallet) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.Wallet).Error
}

func (d *Wallet) GetByWalletID(ctx context.Context, walletID string) (*wallet.Wallet, error) {
	var res *model.Wallet
	err := d.DB.WithContext(ctx).Where("wallet_id = ? ", walletID).First(&res).Error
	return &wallet.Wallet{
		Wallet: res,
	}, err
}

func (d *Wallet) GetWalletsByVaultID(ctx context.Context, vaultID string) ([]*wallet.Wallet, error) {
	var res []*model.Wallet
	err := d.DB.WithContext(ctx).Where("vault_id = ? ", vaultID).Find(&res).Error
	var list []*wallet.Wallet
	for _, each := range res {
		list = append(list, &wallet.Wallet{Wallet: each})
	}
	return list, err
}
