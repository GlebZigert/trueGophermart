# cmd/gophermart

В данной директории будет содержаться код накопительной системы лояльности, который скомпилируется в бинарное
приложение.

Билд автотеста:

go test .\autotests\main_test.go .\autotests\TestFlagRunAddrSuite.go .\autotests\TestEnvRunAddrSuite.go .\autotests\TestRegSuite.go .\autotests\flags.go -c -o=./

Билд целевого сервера:

go build .\cmd\gophermart\main.go 
go1.20.7 build cmd/gophermart/main.go 

Запуск автотеста:

.\autotests.test.exe -binary-path=C:\gophermart\main.exe 

