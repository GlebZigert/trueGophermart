package app

import (
	"fmt"

	"github.com/GlebZigert/gophermart/internal/config"
	"github.com/GlebZigert/gophermart/internal/dblayer"
	"github.com/GlebZigert/gophermart/internal/server"
)

func Run() error {
	config.ParseFlags()
	err := dblayer.Init()
	if err != nil {
		fmt.Println("err: ", err.Error())
		return err
	}
	server.InitRouter()

	return nil
}
