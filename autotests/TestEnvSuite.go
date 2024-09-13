package autotests

import (
	"context"
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/GlebZigert/trueGophermart/internal/fork"

	"github.com/stretchr/testify/suite"
)

/*
Сервис должен поддерживать конфигурирование следующими методами:

    адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
    адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
*/

type TestEnvSuite struct {
	suite.Suite
	serverAddress string
	serverProcess *fork.BackgroundProcess
}

func (suite *TestEnvSuite) SetupSuite() {
	suite.T().Logf("TestEnvSuite SetupSuite")
	suite.Require().NotEmpty(flagTargetBinaryPath, "-binary-path non-empty flag required")
	suite.Require().NotEmpty(flagServerPort, "-server-port non-empty flag required")
	suite.Require().NotEmpty(flagGophermartDatabaseURI, "-gophermart-database-uri non-empty flag required")
	// приравниваем адрес сервера
	suite.serverAddress = "localhost:" + flagServerPort

	// запускаем процесс тестируемого сервера
	{

		envs := append(os.Environ(), []string{
			"RUN_ADDR=" + suite.serverAddress,
			"DATABASE_URI=" + flagGophermartDatabaseURI,
		}...)
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
		port := flagServerPort
		err = p.WaitPort(ctx, "tcp", port)
		if err != nil {
			suite.T().Errorf("Не удалось дождаться пока порт %s станет доступен для запроса: %s", port, err)

			if out := p.Stderr(ctx); len(out) > 0 {
				suite.T().Logf("Получен STDERR лог агента:\n\n%s\n\n", string(out))
			}
			if out := p.Stdout(ctx); len(out) > 0 {
				suite.T().Logf("Получен STDOUT лог агента:\n\n%s\n\n", string(out))
			}

			return
		}

		if out := p.Stderr(ctx); len(out) > 0 {
			suite.T().Logf("Получен STDERR лог агента:\n\n%s\n\n", string(out))
		}
		if out := p.Stdout(ctx); len(out) > 0 {
			suite.T().Logf("Получен STDOUT лог агента:\n\n%s\n\n", string(out))
		}
	}
}

// TearDownSuite высвобождает имеющиеся зависимости
func (suite *TestEnvSuite) TearDownSuite() {
	// посылаем процессу сигналы для остановки
	exitCode, err := suite.serverProcess.Stop(syscall.SIGINT, syscall.SIGKILL)
	if err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return
		}
		suite.T().Logf("Не удалось остановить процесс с помощью сигнала ОС: %s", err)
		return
	}

	// проверяем код завешения
	if exitCode > 0 {
		suite.T().Logf("Процесс завершился с не нулевым статусом %d", exitCode)
	}

	// получаем стандартные выводы (логи) процесса
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	out := suite.serverProcess.Stderr(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDERR лог процесса:\n\n%s", string(out))
	}
	out = suite.serverProcess.Stdout(ctx)
	if len(out) > 0 {
		suite.T().Logf("Получен STDOUT лог процесса:\n\n%s", string(out))
	}
}

func (suite *TestEnvSuite) TestHandlers() {
	// генерируем новый псевдорандомный URL
	//suite.T().Logf("just2")

}
