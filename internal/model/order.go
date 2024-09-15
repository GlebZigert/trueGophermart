package model

type Order struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

type OrderErr struct {
	Err string
}

var FoundNoOrder *OrderErr = &OrderErr{"Не найдено заказов"}

func (e *OrderErr) Error() string {
	return e.Err
}
