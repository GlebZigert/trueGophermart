package autotests

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/GlebZigert/gophermart/internal/fork"
	"github.com/go-resty/resty/v2"

	"github.com/stretchr/testify/suite"
)

/*

Хендлер: POST /api/user/register.
Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
После успешной регистрации должна происходить автоматическая аутентификация пользователя.
Для передачи аутентификационных данных используйте механизм cookies или HTTP-заголовок Authorization.
Формат запроса:

POST /api/user/register HTTP/1.1
Content-Type: application/json
...

{
    "login": "<login>",
    "password": "<password>"
}

Возможные коды ответа:

    200 — пользователь успешно зарегистрирован и аутентифицирован;
    400 — неверный формат запроса;
    409 — логин уже занят;
    500 — внутренняя ошибка сервера.

*/

type TestRegSuite struct {
	suite.Suite
	serverAddress string
	serverProcess *fork.BackgroundProcess
}

func (suite *TestRegSuite) SetupSuite() {
	suite.T().Logf("TestEnvRunAddrSuite SetupSuite")
	suite.Require().NotEmpty(flagTargetBinaryPath, "-binary-path non-empty flag required")
	suite.Require().NotEmpty(flagServerPort, "-server-port non-empty flag required")
	suite.Require().NotEmpty(flagGophermartDatabaseURI, "-gophermart-database-uri non-empty flag required")
	// приравниваем адрес сервера
	suite.serverAddress = "127.0.0.1:" + flagServerPort

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
func (suite *TestRegSuite) TearDownSuite() {
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

func (suite *TestRegSuite) TestHandler() {
	//послать запрос

	// создаем HTTP клиент без поддержки редиректов
	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	httpc := resty.New().
		SetBaseURL("http://" + suite.serverAddress).
		SetRedirectPolicy(redirPolicy)

	suite.Run("shorten", func() {
		// весь тест должен проходить менее чем за 10 секунд
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//шлем запрос  получение списка загруженных пользователем номеров заказов - без ключа авторизации
		suite.T().Logf("Шлю запрос GET orders - без авторизации. Должен прийти ответ со статусом StatusUnauthorized")
		req := httpc.R().
			SetHeader("Content-Type", "application/json").
			SetContext(ctx)
		//я должен получить ответ
		//провожу роверку на наличие ответа
		resp, err := req.Get("/api/user/orders")
		noRespErr := suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}
		//я должен получить ответ со статусом StatusUnauthorized о том что запрос не обработан из за отсутствия валидного ключа авторизации
		////провожу роверку на наличие статуса StatusUnauthorized
		suite.Assert().Equalf(http.StatusUnauthorized, resp.StatusCode(),
			"Несоответствие статус кода ответа ожидаемому в хендлере '%s %s'", req.Method, req.URL)

		// делаем запрос на регистрацию без нормальной посылки - должен быть ответ со статусом 400 - неверный формат запроса
		req = httpc.R().
			SetHeader("Content-Type", "application/json").
			//	SetBody(m).
			SetContext(ctx)

		resp, err = req.Post("/api/user/register")

		//Должны получить ответ со статусом 200 — пользователь успешно зарегистрирован и аутентифицирован;
		//В ответе должен быть HTTP-заголовок Authorization

		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")

		if !noRespErr {
			suite.T().Errorf(err.Error())
		}

		suite.Assert().Equalf(http.StatusBadRequest, resp.StatusCode(),
			"Несоответствие статус кода ответа ожидаемому в хендлере '%s %s'", req.Method, req.URL)

		// делаем запрос на регистрацию c нормальной посылкой - должен быть ответ со статусом 200 и валидным ключем авторизации с userID в хедере

		m := []byte(`{
					"login": "user1",
					"password": "password1"
				}`)

		req = httpc.R().
			SetHeader("Content-Type", "application/json").
			SetBody(m).
			SetContext(ctx)

		resp, err = req.Post("/api/user/register")

		//Должны получить ответ со статусом 200 — пользователь успешно зарегистрирован и аутентифицирован;

		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")

		if !noRespErr {
			suite.T().Errorf(err.Error())
		}

		suite.Assert().Equalf(http.StatusOK, resp.StatusCode(),
			"Несоответствие статус кода ответа ожидаемому в хендлере '%s %s'", req.Method, req.URL)

		//В ответе должен быть HTTP-заголовок Authorization
		authHeader := resp.Header().Get("Authorization")
		setCookieHeader := resp.Header().Get("Set-Cookie")
		suite.Assert().True(authHeader != "" || setCookieHeader != "",
			"Не удалось обнаружить авторизационные данные в ответе")

		// делаем запрос на регистрацию c той же самой посылкой - должен быть ответ со статусом 409 — логин уже занят;

		resp, err = req.Post("/api/user/register")

		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")

		if !noRespErr {
			suite.T().Errorf(err.Error())
		}

		suite.Assert().Equalf(http.StatusConflict, resp.StatusCode(),
			"Несоответствие статус кода ответа ожидаемому в хендлере '%s %s'", req.Method, req.URL)
		//
	})
}
