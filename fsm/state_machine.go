package fsm

import (
	"fmt"
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/logger"
)

type FiniteStateMachine struct {
	cause string // DEPRECATED
	//mode         mode.Mode // DEPRECATED
	logger       logger.Logger
	modelName    string
	states       map[string]*State
	currentState *State
	variables    Variables
	cache        map[string]any
}

func (fsm *FiniteStateMachine) Fire(event string) error {
	fsm.logger.Debugf("Firing %s", event)
	if fsm.currentState == nil {
		//fsm.cause = "No current state"
		//fsm.mode = mode.DEADLOCK
		return fmt.Errorf("fsm_model_fire: no current state")
	}
	fsm.logger.Debugf("Checking %s ...", fsm.currentState.GetName())
	state, err := fsm.currentState.fire(event, &fsm.variables)
	if err != nil {
		return err
	}
	//fsm.mode = currentMode
	if state == nil {
		//fsm.cause = "No resulting state from transition"
		//fsm.mode = mode.DEADLOCK
		return NewDeadlockError("no resulting state from transition")
	}
	newStateName := state
	fsm.logger.Debugf("Transition [%s] -> [%s]", fsm.GetCurrentState().GetName(), *newStateName)
	newState, hasState := fsm.states[*newStateName]
	if !hasState {
		//fsm.cause = "State not found from transition"
		//fsm.mode = mode.CRASH
		//fsm.currentState = nil
		return NewCrashErrorf("state %s not found from transition", *newStateName)
	}
	fsm.currentState = newState
	return nil
}

func (fsm *FiniteStateMachine) GetRegisteredStates() []string {
	cached, contains := fsm.cache["states-keys"]
	if contains {
		return cached.([]string)
	}
	cache := make([]string, 0, len(fsm.states))
	for k := range fsm.states {
		cache = append(cache, k)
	}
	fsm.cache["states-keys"] = cache
	return cache
}

// DEPRECATED
/*func (fsm *FiniteStateMachine) GetMode() mode.Mode {
	return fsm.mode
}*/

// DEPRECATED
func (fsm *FiniteStateMachine) GetCause() string {
	return fsm.cause
}

func (fsm *FiniteStateMachine) GetModelName() string {
	return fsm.modelName
}

func (fsm *FiniteStateMachine) GetCurrentState() *State {
	return fsm.currentState
}

type FsmBuilder struct {
	logger       logger.Logger
	modelName    string
	states       map[string]*State
	initialState *State
	variables    Variables
}

func NewFsmBuilder() FsmBuilder {
	builderLogger := logger.New("FSM (Builder)")
	return FsmBuilder{
		logger:       builderLogger,
		states:       map[string]*State{},
		modelName:    "",
		initialState: nil,
		variables:    NewVariables(),
	}
}

func (fsm *FsmBuilder) Name(modelName string) *FsmBuilder {
	fsm.modelName = modelName
	return fsm
}

func (fsm *FsmBuilder) DeclareVar(key string, value any) *FsmBuilder {
	fsm.variables.Set(key, value)
	return fsm
}

func (fsm *FsmBuilder) Given(state string, f functions.Consumer[*StateBuilder]) *FsmBuilder {
	sb := newStateBuilder(state)
	f(&sb)
	st := sb.build()
	fsm.states[state] = &st
	return fsm
}

func (fsm *FsmBuilder) Initial(state string) *FsmBuilder {
	st, contains := fsm.states[state]
	if !contains {
		//fsm.logger.Error("Initial state '" + state + "' is invalid")
		return fsm
	}
	fsm.initialState = st
	return fsm
}

func (fsm *FsmBuilder) Build() FiniteStateMachine {
	if len(fsm.modelName) == 0 {
		fsm.modelName = "Default (FSM)"
	}
	return FiniteStateMachine{
		//mode:         mode.CONTINUE,
		logger:       logger.New(fsm.modelName),
		modelName:    fsm.modelName,
		currentState: fsm.initialState,
		states:       fsm.states,
		variables:    fsm.variables,
		cache:        map[string]any{},
	}
}
