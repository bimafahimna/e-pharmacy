package cronjob

import (
	"context"
)

type Provider interface {
	Start()
	Stop() context.Context
	AddFunc(crontSchedule string, cmd func()) (EntryID, error)
	AddJob(crontSchedule string, job FuncJob) (EntryID, error)
	Remove(id EntryID)
}

type EntryID int

type FuncJob func()

func (f FuncJob) Run() { f() }
