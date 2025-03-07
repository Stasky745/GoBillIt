package main

import (
	"os"
	"strings"

	"github.com/Stasky745/go-libs/log"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func initalizeLogger() {
	log.InitLogger(setDefaultEnv(APP_PREFIX+"DEBUG", false))
}

func main() {
	initalizeLogger()

	err := app.Run(os.Args)
	log.CheckErr(err, true, "error running app", "cmd", strings.Join(os.Args, " "))
}
