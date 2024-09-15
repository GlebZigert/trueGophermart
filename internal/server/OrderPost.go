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
	"io"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
	"go.uber.org/zap"
)

func (h handler) OrderPost(w http.ResponseWriter, req *http.Request) {
	var err error
	defer packerr.AddErrToReqContext(req, &err)

	//определить что за юзер
	uid, ok := req.Context().Value(config.UIDkey).(int)
	if !ok {
		err = NoUidError

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	//сомтрим на номер заказа
	var number int

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return

	}

	if err = json.Unmarshal(body, &number); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return // err
	}

	logger.Log.Info("Ищу номер заказа  : ", zap.Int("uid", uid), zap.Int("number", number))

	var order model.Order

	result := h.DB.Where("number = ?", number).First(&order)
	if result.Error == nil {
		if order.UID == uid {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			return
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte{})
			return
		}

	}
	order.UID = uid
	order.Number = number

	if result := h.DB.Create(&order); result.Error != nil {

		err = Conflict
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	logger.Log.Info("результат поиска : ", zap.String("err", result.Error.Error()))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte{})
	return
}
