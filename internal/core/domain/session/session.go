package session

import "time"

type Session struct {
	ID                    uint
	UserID                uint
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
