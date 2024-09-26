package model

type User struct {
	ID        int     `gorm:"id"`
	Login     string  `gorm:"login"`
	Password  string  `gorm:"password"`
	Current   float32 `gorm:"current"`
	Withdrawn float32 `gorm:"withdrawn"`
}

type UsersErr struct {
	Err string
}

var FoundNoUser *UsersErr = &UsersErr{"Ffff"}

func (e *UsersErr) Error() string {
	return e.Err
}
