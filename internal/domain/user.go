package domain

import "time"

type User struct {
	ID            uint
	Email         string
	PasswordHash  string
	RegisteredAt  time.Time
	LastVisitedAt time.Time
}
