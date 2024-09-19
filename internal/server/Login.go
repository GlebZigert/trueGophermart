package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

var WrongPassword *model.UsersErr = &model.UsersErr{"Неверный пароль"}

func (srv *Server) Login(w http.ResponseWriter, req *http.Request) {

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

	//если пришла правильно составленная посылка
	srv.logger.Info("try to login: ", map[string]interface{}{

		"login":    user.Login,
		"password": user.Password,
	})
	//проверяем есть ли уже такой логин

	var finded *model.User

	if result := srv.DB.Where("login = ?", user.Login).First(&finded); result.Error != nil {
		//если не нашлось пользователя с таким логином
		err = model.FoundNoUser
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
	jwt, _ := srv.auch.BuildJWTString(finded.ID)
	//добавляю ключ
	w.Header().Add("Authorization", string(jwt))
	//ставлю статус 200
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})

}
