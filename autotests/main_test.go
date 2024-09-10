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

func TestFlagRunAddr(t *testing.T) {
	// Запускает тест-сьют для первой итерации
	suite.Run(t, new(TestFlagRunAddrSuite))
}

func TestEnvRunAddr(t *testing.T) {
	// Запускает тест-сьют для первой итерации
	suite.Run(t, new(TestEnvRunAddrSuite))
}

func TestReg(t *testing.T) {
	// Запускает тест-сьют для первой итерации
	suite.Run(t, new(TestRegSuite))
}
