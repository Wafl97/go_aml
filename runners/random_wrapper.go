package runners

import (
	"math/rand"

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

func RunAsRandom(model *fsm.FiniteStateMachine, iterations int) Summary {
	summary := Summary{
		Path:          make([]string, iterations),
		Occurences:    make(map[string]int, len(model.GetRegisteredStates())),
		DeadlockState: types.None[string](),
	}
	currentState := model.GetCurrentState()
	summary.Path = append(summary.Path, currentState.GetName())
	summary.Occurences[currentState.GetName()] = 1
	log := logger.New("RANDOM WRAPPER")
	log.Infof("Running for %d iterations", iterations)
	for i := 1; i < iterations; i++ {
		//time.Sleep(time.Duration(5) * time.Millisecond)
		arr := currentState.GetEdgeTriggers()
		var currentMode mode.Mode
		if len(arr) > 0 {
			randomChoise := arr[rand.Intn(len(arr))]
			model.Fire(randomChoise)
			currentMode = model.GetMode()
		} else {
			currentMode = mode.DEADLOCK
		}
		switch currentMode {
		case mode.CONTINUE:
			currentState = model.GetCurrentState()
			summary.Path[i] = currentState.GetName()
			summary.Occurences[currentState.GetName()] += 1
		case mode.CRASH:
			log.Errorf("Model crashed. Cause: %s", model.GetCause())
			return summary
		case mode.DEADLOCK:
			log.Error("Deadlock! Exiting ...")
			summary.DeadlockState = types.Some(model.GetCurrentState().GetName())
			return summary
		case mode.TERMINATE:
			log.Infof("Model terminated in state %s", model.GetCurrentState().GetName())
			return summary
		}
	}
	log.Info("Done")
	return summary
}
