package autotests

import (
	"context"
	"os"
	"time"

	"github.com/GlebZigert/gophermart/internal/fork"

	"github.com/stretchr/testify/suite"
)

// Iteration1Suite является сьютом с тестами и состоянием для инкремента
type Iteration1Suite struct {
	suite.Suite
	serverAddress string
	serverProcess *fork.BackgroundProcess
}

func (suite *Iteration1Suite) SetupSuite() {

	suite.Require().NotEmpty(flagTargetBinaryPath, "-binary-path non-empty flag required")

	// прихраниваем адрес сервера
	suite.serverAddress = "http://localhost:8080"

	// запускаем процесс тестируемого сервера
	{
		envs := os.Environ()
		p := fork.NewBackgroundProcess(context.Background(), flagTargetBinaryPath,
			fork.WithEnv(envs...),
		)
		suite.serverProcess = p

		// ожидаем запуска процесса не более 20 секунд
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := p.Start(ctx)
		if err != nil {
			suite.T().Errorf("Невозможно запустить процесс командой %s: %s. Переменные окружения: %+v", p, err, envs)
			return
		}

		// проверяем, что порт успешно занят процессом
		port := "8080"
		err = p.WaitPort(ctx, "tcp", port)
		if err != nil {
			suite.T().Errorf("Не удалось дождаться пока порт %s станет доступен для запроса: %s", port, err)
			return
		}
	}
}

func (suite *Iteration1Suite) TestHandlers() {
	// генерируем новый псевдорандомный URL
	//suite.T().Logf("just2")

}
