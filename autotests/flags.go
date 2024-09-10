package autotests

import "flag"

var (
	flagTargetBinaryPath string // путь до бинарного файла проекта
	flagServerPort       string
)

func init() {
	flag.StringVar(&flagTargetBinaryPath, "binary-path", "", "path to target HTTP server binary")
	flag.StringVar(&flagServerPort, "server-port", "", "port of target address")
}
