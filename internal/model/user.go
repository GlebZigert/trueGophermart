package model

type User struct {
	ID        int     `db:"id"`
	Login     string  `db:"login"`
	Password  string  `db:"password"`
	Current   float32 `db:"current"`
	Withdrawn float32 `db:"withdrawn"`
}

type UsersErr struct {
	Err string
}

var FoundNoUser *UsersErr = &UsersErr{"Ffff"}

func (e *UsersErr) Error() string {
	return e.Err
}
