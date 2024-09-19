package server

import (
	"errors"
	"net/http"

	"github.com/GlebZigert/trueGophermart/internal/config"
	"github.com/GlebZigert/trueGophermart/internal/model"
	"github.com/GlebZigert/trueGophermart/internal/packerr"
)

var NoUidError error = errors.New("Этот реквест прошел проверку в auth но в хэндлере не смог взять uid из контекста")

func (srv *Server) OrderGet(w http.ResponseWriter, req *http.Request) {
	var err error
	defer packerr.AddErrToReqContext(req, &err)
	var orders []model.Order

	//определить что за юзер
	uid, ok := req.Context().Value(config.UIDkey).(int)
	if !ok {
		err = NoUidError

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte{})

	}

	srv.logger.Info("Ищу номера заказов для : ", map[string]interface{}{
		"uid": uid,
	})

	//определить что за юзер
	if result := srv.DB.Find(&orders); result.Error != nil {

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
