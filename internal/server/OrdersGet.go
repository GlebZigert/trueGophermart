package server

import (
	"fmt"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

func (h handler) OrdersGet(w http.ResponseWriter, req *http.Request) {
	var err error
	defer packerr.AddErrToReqContext(req, &err)
	var orders []model.Order

	//определить что за юзер
	if result := h.DB.Find(&orders); result.Error != nil {
		fmt.Println(result.Error)
		err = result.Error

	}

	if len(orders) == 0 {
		err = model.FoundNoOrder
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		w.Write([]byte{})
		return
	}

	//искать в  дб записи заказов от этого юзера

	//если нет записей - 204

	//есть записи - 200

	//что то идет не так - 500

}
