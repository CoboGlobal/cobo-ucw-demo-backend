package transaction

import (
	"runtime/debug"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-kratos/kratos/v2/log"
)

func (u Usecase) ModuleName() string {
	return "BizTransaction"
}

func (u Usecase) RegisterCronTask(c gocron.Scheduler) ([]gocron.Job, error) {
	duration := DefaultQueryPortalDuration
	if u.ucw.SyncTransactionDuration.Seconds > 0 {
		duration = u.ucw.SyncTransactionDuration.AsDuration()
	}

	u.logger.Infof("BizTransaction duration %v", duration.Seconds())
	job1, err := c.NewJob(gocron.DurationJob(duration),
		gocron.NewTask(
			func() {
				defer func() {
					defer func() {
						if p := recover(); p != nil {
							u.logger.Errorf("panic \n%s", string(debug.Stack()))
						}
					}()
				}()
				if err := u.SyncTransactionTask(); err != nil {
					u.logger.Errorf("SyncTransactionTask err %v", err)
				}
			},
		),
		gocron.WithName("AutoUpdateTransactionStatusFromPortal"),
		gocron.WithTags("tss_request", "portal"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Errorf("Biz Transaction RegisterCronTask err: %v", err)
	}
	job2, err := c.NewJob(gocron.DurationJob(time.Second*20),
		gocron.NewTask(
			func() {
				defer func() {
					defer func() {
						if p := recover(); p != nil {
							u.logger.Errorf("panic \n%s", string(debug.Stack()))
						}
					}()
				}()
				if err := u.SubmitCreatedTransactionTask(); err != nil {
					u.logger.Errorf("SubmitCreatedTransactionTask err %v", err)
				}
			},
		),
		gocron.WithName("SubmitCreatedTransactionTask"),
		gocron.WithTags("tss_request", "portal"),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		log.Errorf("Biz Transaction RegisterCronTask err: %v", err)
	}
	return []gocron.Job{job1, job2}, err
}
