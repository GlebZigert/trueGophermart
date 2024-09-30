package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

var Conflict *model.UsersErr = &model.UsersErr{Err: "Конфликт: логин занят"}

func (srv *Server) Register(w http.ResponseWriter, req *http.Request) {

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

	srv.logger.Info("try to register : ", map[string]interface{}{
		"login":    user.Login,
		"password": user.Password,
	})

	//проверяем есть ли уже такой логин

	var finded model.User

	if result := srv.DB.Where("login = ?", user.Login).First(&finded); result.Error == nil {

		//поднять ошибку о конфликте
		err = Conflict
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte{})
		return
	}

	if result := srv.DB.Create(&user); result.Error != nil {

		err = Conflict
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	jwt, _ := srv.auch.BuildJWTString(user.ID)
	//добавляю ключ
	w.Header().Add("Authorization", string(jwt))
	//ставлю статус 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
