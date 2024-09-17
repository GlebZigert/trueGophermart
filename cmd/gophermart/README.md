# cmd/gophermart

В данной директории будет содержаться код накопительной системы лояльности, который скомпилируется в бинарное
приложение.

Билд автотеста:

go test .\autotests\main_test.go .\autotests\TestFlagSuite.go .\autotests\TestEnvSuite.go .\autotests\TestRegSuite.go .\autotests\flags.go -c -o=./

Билд целевого сервера:

go build .\cmd\gophermart\main.go 
go1.20.7 build cmd/gophermart/main.go 

Запуск автотеста:

Windows
.\autotests.test.exe -binary-path=C:\gophermart\main.exe -server-port=8080 -gophermart-database-uri="host=localhost user=postgres password=qwer dbname=testdb sslmode=disable"

Линкус

cd autotests
go1.21.0 test -binary-path=../main -server-port=8080 
-gophermart-database-uri="host=localhost user=gz password=gzpassword dbname=gzbase sslmode=disable"

go test -binary-path=C:\trueGophermart\main.exe -server-port=8080 -gophermart-database-uri="host=localhost user=postgres password=qwer dbname=testdb sslmode=disable"

Запуск сервера
./main.exe -d="host=localhost user=postgres password=qwer dbname=testdb sslmode=disable"

Установка системной переменной в терминале
set DATABASE_URI=host=localhost user=postgres password=qwer dbname=testdb sslmode=disable

Автотест из трушного бинарника

 C:\go-autotests-0.10.14\go-autotests-0.10.14\gophermarttest-windows-amd64.exe -gophermart-binary-path=C:\gophermart\main.exe -gophermart-database-uri="host=localhost user=postgres password=qwer dbname=testdb sslmode=disable" -gophermart-host=localhost -gophermart-port=8081 -accrual-binary-path=.\accrual_windows_amd64 -accrual-database-uri="host=localhost user=postgres password=qwer dbname=testdb sslmode=disable" -accrual-host=localhost -accrual-port=8082