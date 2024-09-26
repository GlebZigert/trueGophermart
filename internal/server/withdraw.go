package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

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
type OrderWithdraw struct {
	Number string  `json:"order"`
	Sum    float32 `json:"sum"`
}

func (srv *Server) Withdraw(w http.ResponseWriter, req *http.Request) {

	var err error
	defer packerr.AddErrToReqContext(req, &err)

	//распаковываем
	body, err := io.ReadAll(req.Body)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	var orderwithdraw OrderWithdraw

	err = json.Unmarshal(body, &orderwithdraw)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	//смотрим кто пользователь

	UserID, ok := req.Context().Value(config.UserIDkey).(int)
	if !ok {
		err = ErrNoUserID

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte{})
		return
	}

	var user model.User

	res := srv.DB.Where("id=?", UserID).First(&user)

	if res.Error != nil {
		err = res.Error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return

	}

	//видим сколько бонусов хочет списать пользователь и сколько у него есть
	srv.logger.Info("Found balance : ", map[string]interface{}{

		"пользователь":    user.ID,
		"хочет списать":   orderwithdraw.Sum,
		"в оплату заказа": orderwithdraw.Number,
		"его бонусы":      user.Current,
	})
	//если  бонусов недостаточно
	if user.Current < orderwithdraw.Sum {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)

		w.Write([]byte{})
		return
	}

	user.Current = user.Current - orderwithdraw.Sum
	user.Withdrawn = user.Withdrawn + orderwithdraw.Sum
	srv.DB.Save(user)

	srv.DB.Create(&model.Withdraw{UserID: user.ID,
		Number:      orderwithdraw.Number,
		Sum:         orderwithdraw.Sum,
		ProcessedAt: time.Now(),
	})

	w.WriteHeader(http.StatusOK)

	w.Write([]byte{})

}
