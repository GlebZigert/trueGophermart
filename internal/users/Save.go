package users

import (
	"github.com/GlebZigert/trueGophermart/internal/dblayer"
	"github.com/GlebZigert/trueGophermart/internal/logger"
	"go.uber.org/zap"
)

func Save(login, password string) (err error) {

	//делаем запрос в базу

	fields := dblayer.Fields{
		"login":    &login,
		"password": &password,
	}

	err = dblayer.Table("users").Save(nil, fields)
	if err != nil {
		logger.Log.Error("Save: ", zap.String("", err.Error()))
		return
	}
	return
}
