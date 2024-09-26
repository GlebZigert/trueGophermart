package server

import (
	"encoding/json"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

/*
Получение информации о выводе средств
Хендлер: GET /api/user/withdrawals.
Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым. Формат даты — RFC3339.
Формат запроса:

GET /api/user/withdrawals HTTP/1.1
Content-Length: 0

Возможные коды ответа:

	  200 — успешная обработка запроса.
	    Формат ответа:

	200 OK HTTP/1.1
	Content-Type: application/json
	...

	[
	    {
	        "order": "2377225624",
	        "sum": 500,
	        "processed_at": "2020-12-09T16:09:57+03:00"
	    }
	]

204 — нет ни одного списания.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func (srv *Server) WithdrawalsGet(w http.ResponseWriter, req *http.Request) {

	var err error
	defer packerr.AddErrToReqContext(req, &err)

	//кто юзер

	//определить что за юзер
	UserID, ok := req.Context().Value(config.UserIDkey).(int)
	if !ok {
		err = ErrNoUserID

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	//смотрим его списания бонусов
	var withdrawals []model.Withdraw

	if result := srv.DB.Where("UserID = ?", UserID).Find(&withdrawals); result.Error != nil {

		err = result.Error
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte{})
		return
	}

	//их нет - 204
	if len(withdrawals) == 0 {
		err = model.FoundNoOrder
		w.WriteHeader(http.StatusNoContent)

		w.Write([]byte{})
		return
	}

	//они есть - 200

	resp, err := json.Marshal(withdrawals)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return //err

	}

	srv.logger.Info("Найдены списания бонусов : ", map[string]interface{}{
		"withdrawals": string(resp),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
