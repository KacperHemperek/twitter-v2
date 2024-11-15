package models

import "time"

type SessionModel struct {
	ID         string    `json:"id" mapstructure:"id"`
	Expiration time.Time `json:"expiration" mapstructure:"expiration"`
	UserID     string    `json:"userId" mapstructure:"userId"`
}

func (m *SessionModel) IsExpired() bool {
	return time.Now().After(m.Expiration)
}
