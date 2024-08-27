package main

import (
	"flag"
	"os"
	"strings"

	"github.com/Wafl97/go_aml/fsm"
	"github.com/Wafl97/go_aml/util/logger"
)

func main() {
	filename := flag.String("file", "model.aml", "")
	logMode := flag.String("log", "warn", "")
	flag.Parse()
	logger.SetLogLevel(*logMode)
	logger.SetLogOut(logger.LogOutToDateFile)
	log := logger.New("MAIN")

	log.Infof("Loading from %s", *filename)

	fileContentsBytes, err := os.ReadFile(*filename)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fileContents := string(fileContentsBytes)
	if strings.Contains(fileContents, "syntax fsm") {
		model, err := fsm.FromString(fileContents)
		if err != nil {
			log.Error(err.Error())
		}
		fsm.Generate(model)
	}
}
