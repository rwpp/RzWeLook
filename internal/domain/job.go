package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id         int64
	Name       string
	Cron       string
	Executor   string
	Cfg        string
	CancelFunc func() error
}

func (j Job) NextTime() time.Time {
	parse := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := parse.Parse(j.Cron)
	return s.Next(time.Now())
}
