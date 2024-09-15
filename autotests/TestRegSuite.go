package autotests

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/GlebZigert/trueGophermart/internal/fork"
	"github.com/go-resty/resty/v2"

	"github.com/stretchr/testify/suite"

	"github.com/adrianbrad/psqldocker"
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
	psqlcontainer *psqldocker.Container
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

	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)
	/*
		// run a new psql docker container.
		var err error
		suite.psqlcontainer, err = psqldocker.NewContainer(
			usr,
			password,
			dbName,
			psqldocker.WithContainerName(containerName),
		)

		if err != nil {
			suite.T().Errorf("Не запустился контейнер с базой")
			return
		}

		// compose the psql dsn.
		dsn := fmt.Sprintf(
			"user=%s "+
				"password=%s "+
				"dbname=%s "+
				"host=localhost "+
				"port=%s "+
				"sslmode=disable",
			usr,
			password,
			dbName,
			suite.psqlcontainer.Port(),
		)
	*/

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

		//надо почистить базу

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
			suite.T().Logf("errors.Is(err, os.ErrProcessDone): %s", err)
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

	suite.Run("Register", func() {
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

		suite.T().Logf("Шлю запрос GET orders - теперь с авторизацией. Должен прийти ответ со статусом StatusOk")
		req = httpc.R().
			SetHeader("Authorization", authHeader).
			SetContext(ctx)
		// я должен получить ответ
		// провожу роверку на наличие ответа
		resp, err = req.Get("/api/user/orders")
		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}
		// я должен получить ответ со статусом StatusUnauthorized о том что запрос не обработан из за отсутствия валидного ключа авторизации
		// //провожу роверку на наличие статуса StatusUnauthorized

		StatusUnauthorized := suite.Assert().NotEqualf(http.StatusUnauthorized, resp.StatusCode(), "")
		if !StatusUnauthorized {
			suite.T().Fatalf("Неавторизован")
		}

		suite.T().Logf("Шлю запрос на авторизацию - с пустым паролем-логином. Должен прийти ответ со статусом 400 - Неверный формат запроса")
		req = httpc.R().
			SetContext(ctx)
		// я должен получить ответ
		// провожу роверку на наличие ответа
		resp, err = req.Post("/api/user/login")
		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}
		// я должен получить ответ со статусом StatusUnauthorized о том что запрос не обработан из за отсутствия валидного ключа авторизации
		// //провожу роверку на наличие статуса StatusUnauthorized

		suite.Assert().Equalf(http.StatusBadRequest, resp.StatusCode(), "")

		suite.T().Logf("Шлю запрос на авторизацию - с неверным паролем-логином. Должен прийти ответ со статусом 401 ")

		wrong := []byte(`{
		"login": "wrong_user1",
		"password": "wrong_password1"
	}`)
		req = httpc.R().
			SetBody(wrong).
			SetContext(ctx)
		// я должен получить ответ
		// провожу роверку на наличие ответа
		resp, err = req.Post("/api/user/login")
		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}
		// я должен получить ответ со статусом StatusUnauthorized о том что запрос не обработан из за отсутствия валидного ключа авторизации
		// //провожу роверку на наличие статуса StatusUnauthorized

		suite.Assert().Equalf(http.StatusUnauthorized, resp.StatusCode(), "")

		suite.T().Logf("Шлю запрос на авторизацию - с правильным паролем и неверным логином. Должен прийти ответ со статусом 401 ")

		wrong = []byte(`{
		"login": "user1",
		"password": "wrong_password1"
	}`)
		req = httpc.R().
			SetBody(wrong).
			SetContext(ctx)
		// я должен получить ответ
		// провожу роверку на наличие ответа
		resp, err = req.Post("/api/user/login")
		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}
		// я должен получить ответ со статусом StatusUnauthorized о том что запрос не обработан из за отсутствия валидного ключа авторизации
		// //провожу роверку на наличие статуса StatusUnauthorized

		suite.Assert().Equalf(http.StatusUnauthorized, resp.StatusCode(), "")

		suite.T().Logf("Шлю запрос на авторизацию - с правильным паролем и  логином. Должен прийти ответ со статусом 200 ")

		wrong = []byte(`{
		"login": "user1",
		"password": "password1"
	}`)
		req = httpc.R().
			SetBody(wrong).
			SetContext(ctx)
		// я должен получить ответ
		// провожу роверку на наличие ответа
		resp, err = req.Post("/api/user/login")
		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}
		// я должен получить ответ со статусом StatusUnauthorized о том что запрос не обработан из за отсутствия валидного ключа авторизации
		// //провожу роверку на наличие статуса StatusUnauthorized

		suite.Assert().Equalf(http.StatusOK, resp.StatusCode(), "")

		authHeader1 := resp.Header().Get("Authorization")
		suite.Assert().True(authHeader1 != "")

	})

	suite.Run("Order", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		suite.T().Logf("Шлю запрос на авторизацию - с правильным паролем и  логином. Должен прийти ответ со статусом 200 ")

		lp := []byte(`{
		"login": "user1",
		"password": "password1"
	}`)
		req := httpc.R().
			SetBody(lp).
			SetContext(ctx)

		resp, err := req.Post("/api/user/login")
		noRespErr := suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}

		suite.Assert().Equalf(http.StatusOK, resp.StatusCode(), "")

		authHeader := resp.Header().Get("Authorization")
		suite.Assert().True(authHeader != "")

		suite.T().Logf("Шлю запрос GET orders - c авторизацией. Должен прийти ответ со статусом StatusUnauthorized")
		req = httpc.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", authHeader).
			SetContext(ctx)

		resp, err = req.Get("/api/user/orders")
		noRespErr = suite.Assert().NoError(err, "Ошибка при попытке сделать запрос")
		if !noRespErr {
			suite.T().Errorf(err.Error())
		}

		suite.Assert().Equalf(http.StatusNoContent, resp.StatusCode(),
			"Несоответствие статус кода ответа ожидаемому в хендлере '%s %s'", req.Method, req.URL)

	})
}
