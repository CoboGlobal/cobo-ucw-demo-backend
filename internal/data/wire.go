package data

import (
	"cobo-ucw-backend/internal/biz/transaction"
	"cobo-ucw-backend/internal/biz/user"
	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/biz/wallet"
	"cobo-ucw-backend/internal/data/dao"
	"cobo-ucw-backend/internal/data/database"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(user.Repo), new(*dao.User)), dao.NewUser,
	wire.Bind(new(user.UserVaultRepo), new(*dao.UserVault)), dao.NewUserVault,
	wire.Bind(new(user.UserNodeRepo), new(*dao.UserNode)), dao.NewUserNode,
	wire.Bind(new(vault.Repo), new(*dao.Vault)), dao.NewVault,
	wire.Bind(new(vault.GroupRepo), new(*dao.Group)), dao.NewGroup,
	wire.Bind(new(vault.GroupNodeRepo), new(*dao.GroupNode)), dao.NewGroupNode,
	wire.Bind(new(vault.TssRequestRepo), new(*dao.TssRequest)), dao.NewTssRequest,
	wire.Bind(new(wallet.Repo), new(*dao.Wallet)), dao.NewWallet,
	wire.Bind(new(transaction.Repo), new(*dao.Transaction)), dao.NewTransaction,
	wire.Bind(new(wallet.AddressRepo), new(*dao.Address)), dao.NewAddress,
	database.NewData,
)
