package test1

import "time"

type ClipboardEntry struct {
	ID        uint `gorm:"primaryKey"`
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Timestamp int
	VideoID   string
	Notes     string
}
