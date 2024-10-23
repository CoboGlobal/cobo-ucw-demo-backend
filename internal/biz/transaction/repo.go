package transaction

import (
	"context"

	v1 "cobo-ucw-backend/api/ucw/v1"
)

type Repo interface {
	Update(ctx context.Context, tx Transaction) error
	Save(ctx context.Context, tx *Transaction) (int64, error)
	GetByTransactionID(ctx context.Context, transactionID string) (*Transaction, error)
	ListTransactions(ctx context.Context, params ListTransactionParams) ([]*Transaction, error)
	HardDelete(ctx context.Context, data *Transaction) error
}

type ListTransactionParams struct {
	WalletID        string
	TokenID         string
	TransactionType v1.Transaction_Type
	Limit           int
	LastID          int64
	Status          []v1.Transaction_Status
	ExternalID      string
}
