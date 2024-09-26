package model

import "time"

type Withdraw struct {
	UID         int       `json :"-",db:"uid"`
	Number      string    `json :"order",db:"number"`
	Sum         float32   `json :"sum",db:"sum"`
	ProcessedAt time.Time `json :"processed_at",db:"processed_at"`
}
