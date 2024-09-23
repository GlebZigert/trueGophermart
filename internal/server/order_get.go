package server

import (
	"encoding/json"
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

	//искать в  дб записи заказов от этого юзера
	srv.logger.Info("Ищу номера заказов для : ", map[string]interface{}{
		"uid": uid,
	})

	if result := srv.DB.Where("UID = ?", uid).Find(&orders); result.Error != nil {

		err = result.Error
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte{})
		return
	}

	//если нет записей - 204

	if len(orders) == 0 {
		err = model.FoundNoOrder
		w.WriteHeader(http.StatusNoContent)

		w.Write([]byte{})
		return
	}

	resp, err := json.Marshal(orders)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return //err

	}

	srv.logger.Info("Найдены заказы : ", map[string]interface{}{
		"orders": string(resp),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	//что то идет не так - 500

}
