package model

import "time"

type Withdraw struct {
	UID         int       `json:"-"`
	Number      string    `json:"order" `
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
