package wallet

import (
	"context"
)

type Repo interface {
	Update(ctx context.Context, wallet Wallet) error
	Save(ctx context.Context, wallet *Wallet) (int64, error)
	GetByWalletID(ctx context.Context, walletID string) (*Wallet, error)
	GetWalletsByVaultID(ctx context.Context, vaultID string) ([]*Wallet, error)
}

type AddressRepo interface {
	Update(ctx context.Context, data Address) error
	Save(ctx context.Context, data *Address) (int64, error)
	ListAddress(ctx context.Context, walletID, chainID string) ([]*Address, error)
}
