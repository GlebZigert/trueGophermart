package server

import "net/http"

func (h handler) OrdersGet(w http.ResponseWriter, req *http.Request) {

	//определить что за юзер

	//искать в  дб записи заказов от этого юзера

	//если нет записей - 204

	//есть записи - 200

	//что то идет не так - 500

}
