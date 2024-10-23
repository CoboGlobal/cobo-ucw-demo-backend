package biz

import (
	"cobo-ucw-backend/internal/biz/transaction"
	"cobo-ucw-backend/internal/biz/vault"
	"cobo-ucw-backend/internal/cron"

	"github.com/google/wire"
)

var CronProviderSet = wire.NewSet(
	NewCronRegisters,
	cron.NewCron)

func NewCronRegisters(
	usecase *transaction.Usecase,
	vaultUsecase *vault.Usecase,
) []cron.Register {
	return []cron.Register{usecase, vaultUsecase}
}
