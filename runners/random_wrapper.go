package runners

import (
	"math/rand"
	"strconv"

	"github.com/Wafl97/go_aml/fsm"
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
		randomChoise := arr[rand.Intn(len(arr))]
		if !fsm.Fire(randomChoise) {
			log.Error("Deadlock! Exiting ...")
			summary.DeadlockState = types.Some(fsm.GetCurrentState().Get().GetName())
			break
		}
		currentState = fsm.GetCurrentState().Get()
		summary.Path = append(summary.Path, currentState.GetName())
		summary.Occurences[currentState.GetName()] += 1
	}
	log.Info("Done")
	return summary
}
