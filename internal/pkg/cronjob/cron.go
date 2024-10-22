package cronjob

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/robfig/cron/v3"
)

type cronProvider struct {
	daemon *cron.Cron
}

func NewCronProvider(config config.CronConfig) Provider {
	return &cronProvider{
		daemon: cron.New(cron.WithLocation(config.TimeLocation)),
	}
}

func (p *cronProvider) Start() {
	logger.Log.Info("Starting Cron...")
	p.daemon.Start()
	logger.Log.Info("Cron is now running")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	p.Stop()
	logger.Log.Info("cron shutdown gracefully")
}

func (p *cronProvider) Stop() context.Context {
	logger.Log.Info("Shutting down cron...")
	return p.daemon.Stop()
}

func (p *cronProvider) AddFunc(cronSchedule string, cmd func()) (EntryID, error) {
	entryID, err := p.daemon.AddFunc(cronSchedule, cmd)
	logger.Log.Infof("Added new cron job with entry id: %d", entryID)
	return EntryID(entryID), err
}

func (p *cronProvider) AddJob(cronSchedule string, job FuncJob) (EntryID, error) {
	entryID, err := p.daemon.AddJob(cronSchedule, job)
	logger.Log.Infof("Added new cron job with entry id: %d", entryID)
	return EntryID(entryID), err
}

func (p *cronProvider) Remove(id EntryID) {
	p.daemon.Remove(cron.EntryID(id))
	logger.Log.Infof("Removed cron job with id: %d", id)
}
