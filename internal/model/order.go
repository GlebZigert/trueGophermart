package model

const ORDER_REGISTERED string = "REGISTERED"

type Order struct {
	ID      int    `db:"id"`
	UID     int    `db:"uid"`
	Number  int    `db:"number"`
	Accrual int    `db:"aqrual"`
	Status  string `db:"status"`
}

type Answer struct {
	Number  string  `json:"order"`
	Accrual float32 `json:"accrual"`
	Status  string  `json:"status"`
}

type OrderErr struct {
	Err string
}

var FoundNoOrder *OrderErr = &OrderErr{"Не найдено заказов"}

func (e *OrderErr) Error() string {
	return e.Err
}
