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

func TestIteration1(t *testing.T) {
	// Запускает тест-сьют для первой итерации
	suite.Run(t, new(Iteration1Suite))
}
