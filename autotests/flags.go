package autotests

import "flag"

var (
	flagTargetBinaryPath      string // путь до бинарного файла проекта
	flagServerPort            string
	flagGophermartDatabaseURI string

	flagAccrualBinaryPath  string
	flagAccrualHost        string
	flagAccrualPort        string
	flagAccrualDatabaseURI string
)

func init() {
	flag.StringVar(&flagTargetBinaryPath, "binary-path", "", "path to target HTTP server binary")
	flag.StringVar(&flagServerPort, "server-port", "", "port of target address")
	flag.StringVar(&flagGophermartDatabaseURI, "gophermart-database-uri", "", "connection string to gophermart database")

	flag.StringVar(&flagAccrualBinaryPath, "accrual-binary-path", "", "path to accrual HTTP server binary")
	flag.StringVar(&flagAccrualHost, "accrual-host", "", "host to run accrual HTTP server on")
	flag.StringVar(&flagAccrualPort, "accrual-port", "", "port to run accrual HTTP server on")
	flag.StringVar(&flagAccrualDatabaseURI, "accrual-database-uri", "", "connection string to accrual database")
}
