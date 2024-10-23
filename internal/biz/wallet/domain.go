package wallet

import (
	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/internal/data/model"
)

type Wallet struct {
	*model.Wallet
}

func (w *Wallet) GetWalletID() string {
	if w != nil {
		if w.Wallet != nil {
			return w.Wallet.WalletID
		}
	}
	return ""
}

type Token struct {
	*v1.Token
}

type TokenBalance struct {
	Balance *v1.TokenBalance
}

func (t *TokenBalance) ToProto() *v1.TokenBalance {
	if t == nil {
		return nil
	}
	return t.Balance
}

type WalletToken struct {
	*v1.Wallet
	Token *v1.TokenAddresses
}

type Address struct {
	*model.Address
}

func (a *Address) ToProto() *v1.Address {
	if a == nil || a.Address == nil {
		return &v1.Address{}
	}
	return &v1.Address{
		Address:  a.Address.Address,
		ChainId:  a.ChainID,
		WalletId: a.WalletID,
		Path:     a.Path,
		Pubkey:   a.PubKey,
		Encoding: a.Encoding,
	}
}

func (t *Token) ToProto() *v1.Token {
	if t == nil {
		return nil
	}
	return t.Token
}

type TokenBalances []*v1.TokenBalance

func (w *Wallet) ToProto() *v1.Wallet {
	if w == nil || w.Wallet == nil {
		return nil
	}
	return &v1.Wallet{
		WalletId: w.WalletID,
		Name:     w.Name,
	}
}

func (w *Wallet) ToProtoWalletInfo() *v1.WalletInfo {
	return &v1.WalletInfo{
		Wallet: w.ToProto(),
	}
}
