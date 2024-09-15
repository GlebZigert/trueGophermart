package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/auth"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
	"go.uber.org/zap"
)

var Conflict *model.UsersErr = &model.UsersErr{"Конфликт: логин занят"}

func (h handler) Register(w http.ResponseWriter, req *http.Request) {
	logger.Log.Info("register-->")

	var err error
	defer packerr.AddErrToReqContext(req, &err)

	var user model.User

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

	var finded model.User

	if result := h.DB.Where("login = ?", user.Login).First(&finded); result.Error == nil {

		//поднять ошибку о конфликте
		err = Conflict
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte{})
		return
	}

	if result := h.DB.Create(&user); result.Error != nil {

		err = Conflict
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	jwt, _ := auth.BuildJWTString(user.ID)
	//добавляю ключ
	w.Header().Add("Authorization", string(jwt))
	//ставлю статус 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
