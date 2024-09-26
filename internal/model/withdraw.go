package model

import "time"

type Withdraw struct {
	UID         int       `db:"uid",json "-"`
	Number      int       `db:"number",json "order"`
	Sum         float32   `db:"sum",json "sum"`
	ProcessedAt time.Time `db:"processed_at",json "processed_at"`
}
