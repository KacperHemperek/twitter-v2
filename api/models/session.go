package models

import "time"

type SessionModel struct {
	ID         string    `json:"id" dbmap:"id"`
	Expiration time.Time `json:"expiration" dbmap:"expiration"`
	UserID     string    `json:"userId" dbmap:"userId"`
}

func (m *SessionModel) IsExpired() bool {
	return time.Now().After(m.Expiration)
}
