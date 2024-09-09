package app

import (
	"github.com/GlebZigert/gophermart/internal/server"
)

func Run() error {

	server.InitRouter()

	return nil
}
