package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/gophermart/internal/auth"
	"github.com/GlebZigert/gophermart/internal/logger"
	"github.com/GlebZigert/gophermart/internal/packerr"
	"github.com/GlebZigert/gophermart/internal/users"
	"go.uber.org/zap"
)

var Conflict *users.UsersErr = &users.UsersErr{"Конфликт: логин занят"}

func Register(w http.ResponseWriter, req *http.Request) {
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

	//если пришла правильная посылка
	logger.Log.Info("try to register: ", zap.String("login", user.Login), zap.String("password", user.Password))
	//проверяем есть ли уже такой логин

	if _, ok := users.Find(user.Login); ok == nil {

		//поднять ошибку о конфликте
		err = Conflict
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte{})
		return
	}

	jwt, _ := auth.BuildJWTString()
	//добавляю ключ
	w.Header().Add("Authorization", string(jwt))
	//ставлю статус 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
