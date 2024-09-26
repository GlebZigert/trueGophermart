package model

import "time"

type Withdraw struct {
	Uid         int       `json:"-"`
	Number      string    `json:"order" `
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
