package autotests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

//go:generate go test -c -o=./autotest

func TestMain(m *testing.M) {
	// Основной тест, запускает все остальные тесты

	os.Exit(m.Run())
}

/*
Сервис должен поддерживать конфигурирование следующими методами:

	адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
*/

func TestFlagRunAddr(t *testing.T) {
	// Запускает тест-сьют для первой итерации
	suite.Run(t, new(TestFlagRunAddrSuite))
}

func TestEnvRunAddr(t *testing.T) {
	// Запускает тест-сьют для первой итерации
	suite.Run(t, new(TestEnvRunAddrSuite))
}
