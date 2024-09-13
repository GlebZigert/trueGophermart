package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/auth"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
	"github.com/GlebZigert/trueGophermart/internal/users"
	"go.uber.org/zap"
)

var WrongPassword *users.UsersErr = &users.UsersErr{"Неверный пароль"}

func Login(w http.ResponseWriter, req *http.Request) {

	var err error
	defer packerr.AddErrToReqContext(req, &err)

	var user users.User

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return //err
	}

	if err = json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return // err
	}

	//если пришла правильно составленная посылка
	logger.Log.Info("try to logib: ", zap.String("login", user.Login), zap.String("password", user.Password))
	//проверяем есть ли уже такой логин

	finded, ok := users.Find(user.Login)
	if ok != nil {

		//если не нашлось пользователя с таким логином
		err = users.FoundNoUser
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte{})
		return
	}

	//если пользователь с таким логином найден - проверяем пароль

	if finded.Password != user.Password {
		//если пароли несовпадают
		err = WrongPassword
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte{})
		return
	}

	//если пароли совпадают - авторизовываем - даем ключ
	jwt, _ := auth.BuildJWTString()
	//добавляю ключ
	w.Header().Add("Authorization", string(jwt))
	//ставлю статус 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})

}
