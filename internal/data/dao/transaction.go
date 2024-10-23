package dao

import (
	"context"
	"strings"

	"cobo-ucw-backend/internal/biz/transaction"
	"cobo-ucw-backend/internal/data/database"
	"cobo-ucw-backend/internal/data/model"

	"gorm.io/gorm/clause"
)

func NewTransaction(db *database.Data) *Transaction {
	return &Transaction{db}
}

type Transaction struct {
	*database.Data
}

func (d *Transaction) Update(ctx context.Context, tx transaction.Transaction) error {
	return d.DB.WithContext(ctx).Model(&tx.Transaction).Updates(tx.Transaction).Error
}

func (d *Transaction) Save(ctx context.Context, tx *transaction.Transaction) (int64, error) {
	return int64(tx.ID), d.DB.WithContext(ctx).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(tx.Transaction).Error
}

func (d *Transaction) GetByTransactionID(ctx context.Context, transactionID string) (*transaction.Transaction, error) {
	var res *model.Transaction
	err := d.DB.WithContext(ctx).Where("transaction_id = ? ", transactionID).First(&res).Error
	return &transaction.Transaction{
		Transaction: res,
	}, err
}

func (d *Transaction) ListTransactions(ctx context.Context, params transaction.ListTransactionParams) ([]*transaction.Transaction, error) {
	var res []*model.Transaction

	var query []string
	var args []interface{}

	if params.WalletID != "" {
		query = append(query, " wallet_id = ? ")
		args = append(args, params.WalletID)
	}

	if params.TransactionType != 0 {
		query = append(query, "type = ? ")
		args = append(args, params.TransactionType)
	}

	if params.TokenID != "" {
		query = append(query, " token_id = ? ")
		args = append(args, params.TokenID)
	}

	if len(params.Status) != 0 {
		query = append(query, " status in (?) ")
		args = append(args, params.Status)
	}

	if params.LastID != 0 {
		query = append(query, " id > ? ")
		args = append(args, params.LastID)
	}

	if params.ExternalID != "" {
		query = append(query, " external_id = ? ")
		args = append(args, params.ExternalID)
	}

	prepare := d.DB.WithContext(ctx).Where(strings.Join(query, " and "), args...)
	var err error
	if params.Limit != 0 {
		err = prepare.Limit(params.Limit).Find(&res).Error
	} else {
		err = prepare.Find(&res).Error
	}
	if err != nil {
		return nil, err
	}

	var list []*transaction.Transaction

	for _, each := range res {
		list = append(list, &transaction.Transaction{Transaction: each})
	}
	return list, err
}

func (d *Transaction) HardDelete(ctx context.Context, data *transaction.Transaction) error {
	return d.DB.WithContext(ctx).Unscoped().Delete(data.Transaction).Error
}
