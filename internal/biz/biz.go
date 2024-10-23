package biz

import (
	"context"

	v1 "cobo-ucw-backend/api/ucw/v1"
	"cobo-ucw-backend/integration/portal"
	"cobo-ucw-backend/internal/biz/transaction"
	"cobo-ucw-backend/internal/biz/user"
	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/biz/wallet"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	wire.Bind(new(TransactionUsecase), new(*transaction.Usecase)), transaction.NewUsecase,
	wire.Bind(new(VaultUsecase), new(*vault.Usecase)), vault.NewUsecase,
	wire.Bind(new(UserUsecase), new(*user.Usecase)), user.NewUsecase,
	wire.Bind(new(WalletUsecase), new(*wallet.Usecase)), wallet.NewUsecase,
	portal.NewClient,
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, transaction *transaction.Transaction) (string, error)
	PrepareTransaction(ctx context.Context, transaction *transaction.Transaction) (*transaction.EstimatedFee, error)
	ListTransaction(ctx context.Context, params transaction.ListTransactionParams) ([]*transaction.Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*transaction.Transaction, error)
	SyncTransactions(ctx context.Context, transactions []*transaction.Transaction) error
	SyncDepositTransaction(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error)
}

type WalletUsecase interface {
	GetWalletInfo(ctx context.Context, walletID string) (*wallet.Wallet, error)
	GetWalletToken(ctx context.Context, walletID string, tokenID string) (*wallet.WalletToken, error)
	ListWalletTokens(ctx context.Context, walletID string) (wallet.TokenBalances, error)
	AddWalletTokenAddress(ctx context.Context, walletID string, chainID string) (wallet.Address, error)
	GetBalance(ctx context.Context, walletID, tokenID, address string) (*wallet.TokenBalance, error)
	CreateWallet(ctx context.Context, wallet *wallet.Wallet) (string, error)
	GetWalletByVaultID(ctx context.Context, vaultID string) (*wallet.Wallet, error)
	GetWalletAddress(ctx context.Context, walletID string, chainID string) ([]*wallet.Address, error)
}

type VaultUsecase interface {
	CreateVault(ctx context.Context, projectID string) (*vault.Vault, error)
	CreateKeyGroup(ctx context.Context, vaultID string, groupNodes vault.GroupNodes, groupType v1.Group_GroupType) (string, error)
	KeyGen(ctx context.Context, userID, vaultID, sourceGroupID, targetGroupID string, groupType v1.Group_GroupType) (string, error)
	KeyRecover(ctx context.Context, userID, vaultID, sourceGroupID string, targetGroupID string) (string, error)
	GetVault(ctx context.Context, vaultID string) (*vault.Vault, error)
	ListTssRequests(ctx context.Context, userID, nodeID string, status v1.TssRequest_Status) ([]*vault.TssRequest, error)
	ListGroups(ctx context.Context, vaultID, nodeID string, groupType v1.Group_GroupType) (vault.Groups, error)
	GetGroup(ctx context.Context, vaultID, groupID string) (*vault.Group, error)
	GetTssRequest(ctx context.Context, tssRequestID string) (*vault.TssRequest, error)
	SyncTssRequests(ctx context.Context, tssRequests []*vault.TssRequest) error
}

type UserUsecase interface {
	Login(ctx context.Context, email string) (*user.User, error)
	GetUserNodes(ctx context.Context, userID string) (user.Nodes, error)
	BindUserNode(ctx context.Context, userID, nodeID string) (*user.UserNode, error)
	GetUserInfo(ctx context.Context, userID string) (*user.User, error)
	GetUserNodeByNodeIDs(ctx context.Context, nodeIDs []string) ([]*user.UserNode, error)
	BindUserVault(ctx context.Context, userID, vaultID string) error
}
