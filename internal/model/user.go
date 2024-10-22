package model

import "time"

type User struct {
	ID           int64
	Role         string
	Email        string
	PasswordHash *string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type CustomerDetail struct {
	UserID          int64
	ProfilePhotoURL *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
