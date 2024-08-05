package main

import (
	"flag"
	"os"
	"strings"

	"github.com/Wafl97/go_aml/fsm"
	"github.com/Wafl97/go_aml/runners"
	"github.com/Wafl97/go_aml/util/logger"
)

func main() {
	filename := flag.String("file", "model.aml", "")
	logMode := flag.String("log", "warn", "")
	flag.Parse()
	logger.SetLogLevelByString(*logMode)
	log := logger.New("MAIN")

	log.Infof("Loading from %s", *filename)

	fileContentsBytes, err := os.ReadFile(*filename)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fileContents := string(fileContentsBytes)
	if strings.Contains(fileContents, "syntax fsm") {
		fsm.FromString(fileContents).HasValue(func(model fsm.FinitStateMachine) {
			log.Debugf("%v", model)
			summary := runners.RunAsRandom(&model, 100)
			summary.DeadlockState.HasValue(func(s string) {
				log.Errorf("Model reached a deadlock in state %s", s)
			})
			log.Info("Done!")
		}).Else(func() {
			log.Error("Model is invalid ... exiting")
		})
	}

	//runners.RunAsCli(&tm)
}
