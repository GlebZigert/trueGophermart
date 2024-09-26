package server

/*
Загрузка номера заказа
Хендлер: POST /api/user/orders.
Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
Номер заказа может быть проверен на корректность ввода с помощью алгоритма Луна.
Формат запроса:

POST /api/user/orders HTTP/1.1
Content-Type: text/plain
...

12345678903

Возможные коды ответа:

    200 — номер заказа уже был загружен этим пользователем;
    202 — новый номер заказа принят в обработку;
    400 — неверный формат запроса;
    401 — пользователь не аутентифицирован;
    409 — номер заказа уже был загружен другим пользователем;
    422 — неверный формат номера заказа;
    500 — внутренняя ошибка сервера.
*/

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
	"github.com/theplant/luhn"
)

var errBadOrder error = errors.New("ордер не прошел проверку по алгоритму Луна")

func (srv *Server) OrderPost(w http.ResponseWriter, req *http.Request) {
	var err error
	defer packerr.AddErrToReqContext(req, &err)

	//определить что за юзер
	uid, ok := req.Context().Value(config.UIDkey).(int)
	if !ok {
		err = ErrNoUID

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	//сомтрим на номер заказа

	var numberValue int

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return

	}

	if err = json.Unmarshal(body, &numberValue); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return // err
	}

	if !luhn.Valid(numberValue) {
		err = errBadOrder
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte{})
		return
	}

	number := strconv.Itoa(numberValue)

	srv.logger.Info("Ищу номер заказа : ", map[string]interface{}{
		"uid":    uid,
		"number": number,
	})

	var order model.Order

	result := srv.DB.Where("number = ?", number).First(&order)
	if result.Error == nil {
		if order.UID == uid {

			srv.logger.Info("Заказ уже был принят : ", map[string]interface{}{
				"uid":    uid,
				"number": number,
			})

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			return
		} else {

			srv.logger.Info("Заказ уже был принят от другого пользователя : ", map[string]interface{}{
				"uid":    order.UID,
				"number": order.Number,
			})

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte{})
			return
		}

	}

	order.UID = uid
	order.Number = number
	order.Status = model.ORDER_REGISTERED

	if result := srv.DB.Create(&order); result.Error != nil {

		err = Conflict
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	srv.logger.Info("Заказ принят : ", map[string]interface{}{
		"uid":    uid,
		"number": number,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte{})
}
