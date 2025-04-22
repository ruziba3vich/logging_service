package models

import "time"

type Log struct {
	ID         string `gorm:"type:UUID;default:generateUUIDv4()"`
	Message    string
	EventTime  time.Time
	Level      string
	Service    string
	ReceivedAt time.Time `gorm:"autoCreateTime"`
}
