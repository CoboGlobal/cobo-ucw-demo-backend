package vault

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/go-kratos/kratos/v2/log"
)

func (u Usecase) ModuleName() string {
	return "BizVault"
}

func (u Usecase) RegisterCronTask(c gocron.Scheduler) ([]gocron.Job, error) {
	duration := DefaultQueryPortalDuration
	if u.ucw.SyncTssRequestDuration.Seconds > 0 {
		duration = u.ucw.SyncTssRequestDuration.AsDuration()
	}

	u.logger.Infof("BizVault duration %v", duration.Seconds())
	job, err := c.NewJob(gocron.DurationJob(duration),
		gocron.NewTask(
			u.SyncTssRequestTask,
		),
		gocron.WithName("AutoUpdateTssRequestStatusFromPortal"),
		gocron.WithTags("tss_request", "portal"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Errorf("Biz Tss Request RegisterCronTask err: %v", err)
	}
	return []gocron.Job{job}, err
}
