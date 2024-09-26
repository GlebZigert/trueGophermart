package model

const ORDER_REGISTERED string = "NEW"

type Order struct {
	ID      int     `gorm:"id"`
	Uid     int     `gorm:"uid"`
	Number  string  `gorm:"number"`
	Accrual float32 `gorm:"aqrual"`
	Status  string  `gorm:"status"`
}

type Answer struct {
	Number  string  `json:"order"`
	Accrual float32 `json:"accrual"`
	Status  string  `json:"status"`
}

type OrderErr struct {
	Err string
}

var FoundNoOrder *OrderErr = &OrderErr{Err: "Не найдено заказов"}

func (e *OrderErr) Error() string {
	return e.Err
}
