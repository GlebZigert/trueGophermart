package server

import (
	"net/http"

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

func (srv *Server) Widthraw(w http.ResponseWriter, req *http.Request) {

	var err error
	defer packerr.AddErrToReqContext(req, &err)

}
