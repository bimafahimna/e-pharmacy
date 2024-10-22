package model

import "time"

type Partner struct {
	ID               int64
	Name             string
	LogoUrl          string
	YearFounded      int
	ActiveDays       string
	OperationalStart string
	OperationalStop  string
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
