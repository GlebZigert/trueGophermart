package server

import (
	"encoding/json"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

/*
Получение текущего баланса пользователя
Хендлер: GET /api/user/balance.
Хендлер доступен только авторизованному пользователю.
 В ответе должны содержаться данные о текущей сумме баллов лояльности,
  а также сумме использованных за весь период регистрации баллов.
Формат запроса:

GET /api/user/balance HTTP/1.1
Content-Length: 0

Возможные коды ответа:

    200 — успешная обработка запроса.
      Формат ответа:

  200 OK HTTP/1.1
  Content-Type: application/json
  ...

  {
      "current": 500.5,
      "withdrawn": 42
  }


401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/

func (srv *Server) BalanceGet(w http.ResponseWriter, req *http.Request) {

	var err error
	defer packerr.AddErrToReqContext(req, &err)

	//определить что за юзер
	uid, ok := req.Context().Value(config.UIDkey).(int)
	if !ok {
		err = ErrNoUID

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte{})
		return

	}

	var user model.User

	res := srv.DB.Where("id=?", uid).First(&user)

	if res.Error != nil {
		err = res.Error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return

	}

	var balance model.Balance

	balance.Current = user.Current
	balance.Withdrawn = user.Withdrawn

	srv.logger.Info("Found balance : ", map[string]interface{}{
		"current":   balance.Current,
		"withdrawn": balance.Withdrawn,
	})

	resp, err := json.Marshal(balance)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
