package dao

import (
	"context"
	"strings"

	"cobo-ucw-backend/internal/biz/wallet"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm/clause"
)

func NewAddress(db *database.Data) *Address {
	return &Address{db}
}

type Address struct {
	*database.Data
}

func (d *Address) Save(ctx context.Context, data *wallet.Address) (int64, error) {
	return int64(data.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(data.Address).Error

}

func (d *Address) ListAddress(ctx context.Context, walletID, chainID string) ([]*wallet.Address, error) {

	var res []*model.Address
	var query []string
	var args []interface{}

	if walletID != "" {
		query = append(query, " wallet_id = ? ")
		args = append(args, walletID)
	}

	if chainID != "" {
		query = append(query, " chain_id = ? ")
		args = append(args, chainID)
	}

	err := d.DB.WithContext(ctx).Where(strings.Join(query, " and "), args...).Find(&res).Error
	var list []*wallet.Address

	for _, each := range res {
		list = append(list, &wallet.Address{
			Address: each,
		})
	}
	return list, err
}

func (d *Address) Update(ctx context.Context, data wallet.Address) error {
	return d.DB.WithContext(ctx).Model(&data.Address).Updates(data.Address).Error
}
