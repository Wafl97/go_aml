package runners

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Wafl97/go_aml/fsm"
	"github.com/Wafl97/go_aml/fsm/mode"
	"github.com/Wafl97/go_aml/util/logger"
	"github.com/Wafl97/go_aml/util/types"
)

type Summary struct {
	Path          []string
	Occurences    map[string]int
	DeadlockState types.Option[string]
}

func RunAsRandom(fsm *fsm.FinitStateMachine, iterations int) Summary {
	summary := Summary{
		Path:          make([]string, 0, iterations),
		Occurences:    make(map[string]int, len(fsm.GetRegisteredStates())),
		DeadlockState: types.None[string](),
	}
	currentState := fsm.GetCurrentState().Get()
	summary.Path = append(summary.Path, currentState.GetName())
	summary.Occurences[currentState.GetName()] = 1
	log := logger.New("RANDOM WRAPPER")
	log.Info("Running for " + strconv.Itoa(iterations) + " iterations")
	for i := 1; i < iterations; i++ {
		//time.Sleep(time.Duration(5) * time.Millisecond)
		arr := currentState.GetEdgeTriggers()
		var currentMode mode.Mode
		if len(arr) > 0 {
			randomChoise := arr[rand.Intn(len(arr))]
			fsm.Fire(randomChoise)
			currentMode = fsm.GetMode()
		} else {
			currentMode = mode.DEADLOCK
		}
		switch currentMode {
		case mode.CONTINUE:
			currentState = fsm.GetCurrentState().Get()
			summary.Path = append(summary.Path, currentState.GetName())
			summary.Occurences[currentState.GetName()] += 1
		case mode.CRASH:
			log.Error(fmt.Sprintf("Model crashed. Cause: %s", fsm.GetCause()))
			return summary
		case mode.DEADLOCK:
			log.Error("Deadlock! Exiting ...")
			summary.DeadlockState = types.Some(fsm.GetCurrentState().Get().GetName())
			return summary
		case mode.TERMINATED:
			log.Info(fmt.Sprintf("Model terminated in state %s", fsm.GetCurrentState().Get().GetName()))
			return summary
		}
	}
	log.Info("Done")
	return summary
}
