package users

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

type UsersErr struct {
	Err string
}

func (e *UsersErr) Error() string {
	return e.Err
}
