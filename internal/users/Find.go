package users

import (
	"github.com/GlebZigert/gophermart/internal/dblayer"
	"github.com/GlebZigert/gophermart/internal/logger"
	"go.uber.org/zap"
)

var FoundNoUser *UsersErr = &UsersErr{"пользователь не найден"}

func Find(login string) (*User, error) {

	//делаем запрос в базу

	var wanted User
	fields := dblayer.Fields{
		"id":       &wanted.ID,
		"login":    &wanted.Login,
		"password": &wanted.Password,
	}

	rows, values, err := dblayer.Table("users").Seek("login = ?", login).Get(nil, fields)

	defer rows.Close()

	if err != nil {
		logger.Log.Error("database.PingContext: ", zap.String("", err.Error()))
		return nil, err
	}

	if rows == nil {
		logger.Log.Error("database.PingContext: ", zap.String("", "FoundNoUser"))
		return nil, FoundNoUser
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if nil != err {
			break
		}
		return &wanted, nil
	}

	return nil, nil

}
