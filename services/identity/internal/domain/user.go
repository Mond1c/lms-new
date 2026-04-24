package domain

import "time"

type User struct {
	ID           string
	Email        Email
	DisplayName  string
	PasswordHash PasswordHash
	TelegramID   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
