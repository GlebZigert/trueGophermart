package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

/*
Запрос на списание средств
Хендлер: POST /api/user/balance/withdraw
Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя, в счёт оплаты которого списываются баллы.
Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
Формат запроса:

POST /api/user/balance/withdraw HTTP/1.1
Content-Type: application/json

	{
	    "order": "2377225624",
	    "sum": 751
	}

Здесь order — номер заказа, а sum — сумма баллов к списанию в счёт оплаты.
Возможные коды ответа:

	200 — успешная обработка запроса;
	401 — пользователь не авторизован;
	402 — на счету недостаточно средств;
	422 — неверный номер заказа;
	500 — внутренняя ошибка сервера.
*/
type OrderWidthraw struct {
	Number string  `json:"order"`
	Sum    float32 `json:"sum"`
}

func (srv *Server) Widthraw(w http.ResponseWriter, req *http.Request) {

	var err error
	defer packerr.AddErrToReqContext(req, &err)

	//распаковываем
	body, err := io.ReadAll(req.Body)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})

	}

	var orderwidthraw OrderWidthraw

	err = json.Unmarshal(body, &orderwidthraw)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
	}

	//смотрим кто пользователь

	uid, ok := req.Context().Value(config.UIDkey).(int)
	if !ok {
		err = NoUidError

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte{})

	}

	var user model.User

	res := srv.DB.Where("id=?", uid).First(&user)

	if res.Error != nil {
		err = res.Error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return

	}

	//видим сколько бонусов хочет списать пользователь и сколько у него есть
	srv.logger.Info("Found balance : ", map[string]interface{}{

		"польхователь":    user.ID,
		"хочет списать":   orderwidthraw.Sum,
		"в оплату заказа": orderwidthraw.Number,
		"его бонусы":      user.Current,
	})
	//если  бонусов недостаточно
	if user.Current < orderwidthraw.Sum {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)

		w.Write([]byte{})
		return
	}

	user.Current = user.Current - orderwidthraw.Sum
	user.Widthdrawn = user.Widthdrawn + orderwidthraw.Sum
	srv.DB.Save(user)

	w.WriteHeader(http.StatusOK)

	w.Write([]byte{})
	return
}
