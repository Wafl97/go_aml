package runners

import (
	"fmt"
	"math/rand"

	"github.com/Wafl97/go_aml/fsm"
	"github.com/Wafl97/go_aml/util/logger"
)

type Summary struct {
	Path        []string
	Occurrences map[string]int
}

func RunAsRandom(model *fsm.FiniteStateMachine, iterations int) (Summary, error) {
	summary := Summary{
		Path:        make([]string, iterations),
		Occurrences: make(map[string]int, len(model.GetRegisteredStates())),
	}
	currentState := model.GetCurrentState()
	summary.Path = append(summary.Path, currentState.GetName())
	summary.Occurrences[currentState.GetName()] = 1
	log := logger.New("RANDOM WRAPPER")
	log.Infof("Running for %d iterations", iterations)
	for i := 1; i < iterations; i++ { // First state is already applied, so start from 1
		arr := currentState.GetEdgeTriggers()
		if len(arr) == 0 {
			return summary, fsm.NewDeadlockError("no possible transitions")
		}
		randomChoice := arr[rand.Intn(len(arr))]
		err := model.Fire(randomChoice)
		if err != nil {
			return summary, fmt.Errorf("random_wrapper: %w", err)
		}
		currentState = model.GetCurrentState()
		summary.Path[i] = currentState.GetName()
		summary.Occurrences[currentState.GetName()] += 1
	}
	log.Info("Done")
	return summary, nil
}
