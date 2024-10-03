package model

import "sync"

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

func NewUser(id int) *User {
	return &User{id, "", "", 0.0, 0.0}
}

func (user *User) Lock() {
	//здесь у меня должна быть мапа с ключем айди и валуе мютексом

	//если нет адкветаного айди как обработать?

	//если айди норм но он новый - добавляем его в мапу

	getMutex(user.ID).Lock()
}

func (user *User) Unlock() {
	getMutex(user.ID).Unlock()
}

var mapa map[int]*sync.Mutex

func getMutex(id int) *sync.Mutex {

	if mapa == nil {
		mapa = make(map[int]*sync.Mutex)
	}

	mut, ok := mapa[id]

	if !ok {
		mut = &sync.Mutex{}
		mapa[id] = mut
	}

	return mut

}
