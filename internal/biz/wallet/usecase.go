package wallet

import (
	"context"
	"strings"

	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/integration/portal"
	"cobo-ucw-backend/internal/data/model"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/shopspring/decimal"
)

func NewUsecase(
	client *portal.Client,
	walletRepo Repo,
	addressRepo AddressRepo,
	logger log.Logger,
) *Usecase {
	return &Usecase{
		client:      client,
		walletRepo:  walletRepo,
		addressRepo: addressRepo,
		logger:      log.NewHelper(logger),
	}
}

type Usecase struct {
	client      *portal.Client
	walletRepo  Repo
	addressRepo AddressRepo
	logger      *log.Helper
}

func (u Usecase) GetWalletInfo(ctx context.Context, walletID string) (*Wallet, error) {
	wallet, err := u.walletRepo.GetByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (u Usecase) GetWalletByVaultID(ctx context.Context, vaultID string) (*Wallet, error) {
	wallets, err := u.walletRepo.GetWalletsByVaultID(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	if len(wallets) == 0 {
		return &Wallet{}, nil
	}

	return wallets[0], nil
}

func (u Usecase) GetWalletToken(ctx context.Context, walletID string, tokenID string) (*WalletToken, error) {
	wallet, err := u.walletRepo.GetByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	token, err := u.client.GetToken(u.client.WithContext(ctx), tokenID)
	if err != nil {
		u.logger.Errorf("GetWalletToken GetToken token id %s, err %v", tokenID, err)
		return nil, err
	}
	execute, _, err := u.client.WalletsAPI.ListTokenBalancesForWallet(u.client.WithContext(ctx), walletID).TokenIds(tokenID).Execute()
	if err != nil {
		u.logger.Errorf("GetWalletToken ListTokenBalancesForWallet wallet id %s, token id %s, err %v", walletID, tokenID, err)
		return nil, err
	}

	list, err := u.addressRepo.ListAddress(ctx, walletID, token.GetChainId())
	if err != nil {
		return nil, err
	}
	var addresses []*v1.Address
	for _, each := range list {
		addresses = append(addresses, each.ToProto())
	}
	var (
		available = "0"
		total     = "0"
		locked    = "0"
	)

	if len(execute.GetData()) > 0 {
		tokenBalance := execute.GetData()[0]
		available = tokenBalance.Balance.GetAvailable()
		total = tokenBalance.Balance.GetTotal()
		locked = tokenBalance.Balance.GetLocked()
	}

	return &WalletToken{
		Wallet: wallet.ToProto(),
		Token: &v1.TokenAddresses{
			Token: &v1.TokenBalance{
				Token: &v1.Token{
					TokenId: token.GetTokenId(),
					Name:    token.GetName(),
					Decimal: token.GetDecimal(),
					Symbol:  token.GetSymbol(),
					Chain:   token.GetChainId(),
					IconUrl: token.GetIconUrl(),
				},
				Balance:    total,
				AbsBalance: "",
				Available:  available,
				Locked:     locked,
			},
			Addresses: addresses,
		},
	}, nil
}

func (u Usecase) ListWalletTokens(ctx context.Context, walletID string) (TokenBalances, error) {
	enabledTokens, _, err := u.client.WalletsAPI.ListEnabledTokens(u.client.WithContext(ctx)).Execute()
	if err != nil {
		u.logger.Errorf("ListWalletTokens ListEnabledTokens wallet id %s, err %v", walletID, err)
		return nil, err
	}

	var tokenIDs []string
	for _, each := range enabledTokens.GetData() {
		tokenIDs = append(tokenIDs, each.GetTokenId())
	}

	execute, _, err := u.client.WalletsAPI.ListTokenBalancesForWallet(u.client.WithContext(ctx), walletID).TokenIds(strings.Join(tokenIDs, ",")).Execute()
	if err != nil {
		u.logger.Errorf("GetWalletToken ListTokenBalancesForWallet wallet id %s, token ids %v, err %v", walletID, tokenIDs, err)
		return nil, err
	}
	var res TokenBalances

	balanceSet := make(map[string]CoboWaas2.TokenBalance, 0)
	for _, each := range execute.GetData() {
		balanceSet[each.GetTokenId()] = each
	}

	for _, each := range enabledTokens.GetData() {
		token := &v1.Token{
			TokenId: each.GetTokenId(),
			Name:    each.GetName(),
			Decimal: each.GetDecimal(),
			Symbol:  each.GetSymbol(),
			Chain:   each.GetChainId(),
			IconUrl: each.GetIconUrl(),
		}
		balance, ok := balanceSet[each.GetTokenId()]
		if ok {
			available, err := decimal.NewFromString(balance.Balance.GetAvailable())
			if err != nil {
				return nil, err
			}
			total, err := decimal.NewFromString(balance.Balance.GetTotal())
			if err != nil {
				return nil, err
			}
			locked, err := decimal.NewFromString(balance.Balance.GetLocked())
			if err != nil {
				return nil, err
			}
			res = append(res, &v1.TokenBalance{
				Token:      token,
				Balance:    total.String(),
				AbsBalance: "",
				Available:  available.String(),
				Locked:     locked.String(),
			})
		} else {
			res = append(res, &v1.TokenBalance{
				Token:      token,
				Balance:    "0",
				AbsBalance: "",
				Available:  "0",
				Locked:     "0",
			})
		}
	}
	return res, nil
}

func (u Usecase) AddWalletTokenAddress(ctx context.Context, walletID string, tokenID string) (Address, error) {
	tokenInfo, err := u.client.GetToken(u.client.WithContext(ctx), tokenID)
	if err != nil {
		u.logger.Errorf("AddWalletTokenAddress GetToken token id %s, err %v", tokenID, err)
		return Address{}, err
	}
	execute, _, err := u.client.WalletsAPI.CreateAddress(u.client.WithContext(ctx), walletID).CreateAddressRequest(CoboWaas2.CreateAddressRequest{
		ChainId:  tokenInfo.ChainId,
		Count:    1,
		Encoding: nil,
	}).Execute()
	if err != nil {
		u.logger.Errorf("AddWalletTokenAddress CreateAddress wallet id %s, token id %s, err %v", walletID, tokenID, err)
		return Address{}, err
	}

	address := Address{
		Address: &model.Address{
			WalletID: walletID,
			ChainID:  execute[0].GetChainId(),
			Address:  execute[0].GetAddress(),
			Path:     execute[0].GetPath(),
			PubKey:   execute[0].GetPubkey(),
			Encoding: string(execute[0].GetEncoding()),
		},
	}
	_, err = u.addressRepo.Save(ctx, &address)
	if err != nil {
		return Address{}, err
	}
	return address, nil

}

func (u Usecase) GetBalance(ctx context.Context, walletID, tokenID, address string) (*TokenBalance, error) {
	execute, _, err := u.client.WalletsAPI.ListTokenBalancesForAddress(u.client.WithContext(ctx), walletID, address).TokenIds(tokenID).Execute()
	if err != nil {
		u.logger.Errorf("GetBalance ListTokenBalancesForAddress wallet id %s, token id %s, address %s err %v", walletID, tokenID, address, err)
		return nil, err
	}
	tokenInfo, err := u.client.GetToken(u.client.WithContext(ctx), tokenID)
	if err != nil {
		return nil, err
	}
	token := &v1.Token{
		TokenId: tokenInfo.GetTokenId(),
		Name:    tokenInfo.GetName(),
		Decimal: tokenInfo.GetDecimal(),
		Symbol:  tokenInfo.GetSymbol(),
		Chain:   tokenInfo.GetChainId(),
		IconUrl: tokenInfo.GetIconUrl(),
	}
	if len(execute.GetData()) == 0 {
		return &TokenBalance{
			Balance: &v1.TokenBalance{
				Token:      token,
				Balance:    "0",
				AbsBalance: "",
				Available:  "0",
				Locked:     "0",
			},
		}, nil
	}

	balance := execute.GetData()[0].GetBalance()
	return &TokenBalance{
		Balance: &v1.TokenBalance{
			Token:      token,
			Balance:    balance.GetTotal(),
			AbsBalance: "",
			Available:  balance.GetAvailable(),
			Locked:     balance.GetLocked(),
		},
	}, nil
}

func (u Usecase) CreateWallet(ctx context.Context, wallet *Wallet) (string, error) {

	wallets, err := u.walletRepo.GetWalletsByVaultID(ctx, wallet.VaultID)
	if err != nil {
		return "", err
	}

	if len(wallets) > 0 {
		return wallets[0].WalletID, nil
	}
	execute, _, err := u.client.WalletsAPI.CreateWallet(u.client.WithContext(ctx)).CreateWalletParams(CoboWaas2.CreateWalletParams{
		CreateCustodialWalletParams: nil,
		CreateExchangeWalletParams:  nil,
		CreateMpcWalletParams: &CoboWaas2.CreateMpcWalletParams{
			Name:          wallet.Name,
			WalletType:    CoboWaas2.WALLETTYPE_MPC,
			WalletSubtype: CoboWaas2.WALLETSUBTYPE_USER_CONTROLLED,
			VaultId:       wallet.VaultID,
		},
	}).Execute()
	if err != nil {
		u.logger.Errorf("CreateWallet vault id %s, wallet name %s, err %v", wallet.VaultID, wallet.Name, err)
		return "", err
	}

	wallet.WalletID = execute.MPCWalletInfo.GetWalletId()

	_, err = u.walletRepo.Save(ctx, wallet)
	return wallet.WalletID, err
}

func (u Usecase) GetWalletAddress(ctx context.Context, walletID string, chainID string) ([]*Address, error) {
	list, err := u.addressRepo.ListAddress(ctx, walletID, chainID)
	if err != nil {
		return nil, err
	}
	return list, nil
}
