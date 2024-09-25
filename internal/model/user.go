package model

type User struct {
	ID         int     `db:"id"`
	Login      string  `db:"login"`
	Password   string  `db:"password"`
	Current    float64 `db:"current"`
	Widthdrawn float64 `db:"widthdrawn"`
}

type UsersErr struct {
	Err string
}

var FoundNoUser *UsersErr = &UsersErr{"Ffff"}

func (e *UsersErr) Error() string {
	return e.Err
}
