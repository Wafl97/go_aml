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

func RunAsRandom(fsm *fsm.FiniteStateMachine, iterations int) Summary {
	summary := Summary{
		Path:          make([]string, iterations),
		Occurences:    make(map[string]int, len(fsm.GetRegisteredStates())),
		DeadlockState: types.None[string](),
	}
	currentState := fsm.GetCurrentState().Get()
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
			fsm.Fire(randomChoise)
			currentMode = fsm.GetMode()
		} else {
			currentMode = mode.DEADLOCK
		}
		switch currentMode {
		case mode.CONTINUE:
			currentState = fsm.GetCurrentState().Get()
			summary.Path[i] = currentState.GetName()
			summary.Occurences[currentState.GetName()] += 1
		case mode.CRASH:
			log.Errorf("Model crashed. Cause: %s", fsm.GetCause())
			return summary
		case mode.DEADLOCK:
			log.Error("Deadlock! Exiting ...")
			summary.DeadlockState = types.Some(fsm.GetCurrentState().Get().GetName())
			return summary
		case mode.TERMINATE:
			log.Infof("Model terminated in state %s", fsm.GetCurrentState().Get().GetName())
			return summary
		}
	}
	log.Info("Done")
	return summary
}
