package model

import "time"

type Pharmacist struct {
	UserID            int64
	Email             string
	Name              string
	SipaNumber        string
	WhatsappNumber    string
	YearsOfExperience string
	IsAssigned        bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
