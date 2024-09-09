package autotests

import "flag"

var (
	flagTargetBinaryPath string // путь до бинарного файла проекта
)

func init() {
	flag.StringVar(&flagTargetBinaryPath, "binary-path", "", "path to target HTTP server binary")
}
