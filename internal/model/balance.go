package model

type Balance struct {
	Current   float32 `db:"order"`
	Withdrawn float32 `db:"withdrawn"`
}
