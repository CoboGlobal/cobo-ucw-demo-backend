package cron

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type Register interface {
	ModuleName() string
	RegisterCronTask(c gocron.Scheduler) ([]gocron.Job, error)
}

type Cron struct {
	scheduler gocron.Scheduler
	modules   map[string]Register
	log       *log.Helper
}

func NewCron(logger log.Logger, registers ...Register) *Cron {
	log := log.NewHelper(log.With(logger, "module.name", "cron"))
	scheduler, err := gocron.NewScheduler(gocron.WithGlobalJobOptions(
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					if err != nil {
						log.Warnf("cron: run job: %s name: %s err: %v", jobID, jobName, err)
						return
					}
					log.Debugf("cron: run job: %s name: %s success", jobID, jobName)
				},
			),
		)))
	if err != nil {
		panic(err)
	}
	modules := make(map[string]Register)
	cron := &Cron{
		scheduler: scheduler,
		modules:   modules,
		log:       log,
	}
	for _, register := range registers {
		modules[register.ModuleName()] = register
		if err := cron.RegisterJob(register); err != nil {
			panic(err)
		}
	}
	return cron
}

func (c *Cron) Start() {
	go c.scheduler.Start()
	c.log.Infof("start")
}

func (c *Cron) Stop() error {
	err := c.scheduler.Shutdown()
	if err != nil {
		c.log.Infof("stop failed: %v", err)
		return err
	}
	c.log.Infof("stop")
	return nil
}

func (c *Cron) RegisterJob(register Register) error {

	c.modules[register.ModuleName()] = register
	jobs, err := register.RegisterCronTask(c.scheduler)
	if err != nil {
		log.Errorf("register %s job err: %v", register.ModuleName(), err)
		return err
	}
	for _, job := range jobs {
		c.log.Infof("register %s job: %s id: %s tags: %v success", register.ModuleName(), job.Name(), job.ID(), job.Tags())
	}

	return nil
}
