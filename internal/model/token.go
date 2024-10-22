package model

import "time"

type VerificationToken struct {
	ID        int
	UserID    int64
	Token     string
	UsedAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}

type PasswordResetToken struct {
	ID        int
	UserID    int64
	Token     string
	UsedAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt time.Time
}
