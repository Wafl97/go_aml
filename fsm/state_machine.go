package fsm

import (
	"fmt"

	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/logger"
	"github.com/Wafl97/go_aml/util/types"
)

type FinitStateMachine struct {
	logger       logger.Logger
	modelName    string
	states       map[string]*State
	currentState types.Option[*State]
	variables    Variables
	cache        map[string]any
}

func (fsm *FinitStateMachine) Fire(event string) bool {
	success := true
	fsm.logger.Debug("Firing " + event)
	fsm.currentState.HasValue(func(s *State) {
		fsm.logger.Debug("Checking " + s.GetName() + " ...")
		state, termination := s.fire(event, &fsm.variables)
		if termination == DO {
			success = false
			return
		}
		state.HasValue(func(s string) {
			fsm.logger.Debug(fmt.Sprintf("Transition [%s] -> [%s]", fsm.GetCurrentState().Get().GetName(), s))
			fsm.currentState = types.Some(fsm.states[s])
		}).Else(func() {
			success = false
		})
	})
	return success
}

func (fsm *FinitStateMachine) GetRegisteredStates() []string {
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

func (fsm *FinitStateMachine) GetModelName() string {
	return fsm.modelName
}

func (fsm *FinitStateMachine) GetCurrentState() types.Option[*State] {
	return fsm.currentState
}

type FsmBuilder struct {
	logger       logger.Logger
	modelName    string
	states       map[string]*State
	initialState types.Option[*State]
	variables    Variables
}

func NewFsmBuilder() FsmBuilder {
	builderLogger := logger.New("FSM (Builder)")
	return FsmBuilder{
		logger:       builderLogger,
		states:       map[string]*State{},
		modelName:    "",
		initialState: types.None[*State](),
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
	fsm.initialState = types.Some(st)
	return fsm
}

func (fsm *FsmBuilder) Build() FinitStateMachine {
	if len(fsm.modelName) == 0 {
		fsm.modelName = "Default (FSM)"
	}
	return FinitStateMachine{
		logger:       logger.New(fsm.modelName),
		modelName:    fsm.modelName,
		currentState: fsm.initialState,
		states:       fsm.states,
		variables:    fsm.variables,
		cache:        map[string]any{},
	}
}
