package domain

import "time"

type Session struct {
	RefreshToken string    `json:"refreshToken" db:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expiresAt"`
}