package fsm

import (
	"fmt"
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/logger"
)

type FiniteStateMachine struct {
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
		return fmt.Errorf("fsm_model_fire: no current state")
	}
	fsm.logger.Debugf("Checking %s ...", fsm.currentState.GetName())
	state, err := fsm.currentState.fire(event, &fsm.variables)
	if err != nil {
		return err
	}
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

func (fsm *FiniteStateMachine) GetModelName() string {
	return fsm.modelName
}

func (fsm *FiniteStateMachine) GetCurrentState() *State {
	return fsm.currentState
}

type FiniteStateMachineBuilder struct {
	logger       logger.Logger
	modelName    string
	states       map[string]*State
	initialState *State
	variables    Variables
}

func NewFsmBuilder() FiniteStateMachineBuilder {
	builderLogger := logger.New("FSM (Builder)")
	return FiniteStateMachineBuilder{
		logger:       builderLogger,
		states:       map[string]*State{},
		modelName:    "",
		initialState: nil,
		variables:    NewVariables(),
	}
}

func (builder *FiniteStateMachineBuilder) Name(modelName string) *FiniteStateMachineBuilder {
	builder.modelName = modelName
	return builder
}

func (builder *FiniteStateMachineBuilder) DeclareVar(key string, varType VariableType, value any) *FiniteStateMachineBuilder {
	builder.variables.Set(key, varType, value)
	return builder
}

func (builder *FiniteStateMachineBuilder) Given(state string, f functions.Consumer[*StateBuilder]) *FiniteStateMachineBuilder {
	sb := NewStateBuilder()
	sb.name = state
	f(sb)
	st := sb.build()
	builder.states[state] = &st
	return builder
}

func (builder *FiniteStateMachineBuilder) Given2(state string, stateBuilder *StateBuilder) *FiniteStateMachineBuilder {
	stateBuilder.name = state
	st := stateBuilder.build()
	builder.states[state] = &st
	return builder
}

func (builder *FiniteStateMachineBuilder) Initial(state string) *FiniteStateMachineBuilder {
	st, contains := builder.states[state]
	if !contains {
		//fsm.logger.Error("Initial state '" + state + "' is invalid")
		return builder
	}
	builder.initialState = st
	return builder
}

func (builder *FiniteStateMachineBuilder) Build() FiniteStateMachine {
	if len(builder.modelName) == 0 {
		builder.modelName = "Default (FSM)"
	}
	return FiniteStateMachine{
		logger:       logger.New(builder.modelName),
		modelName:    builder.modelName,
		currentState: builder.initialState,
		states:       builder.states,
		variables:    builder.variables,
		cache:        map[string]any{},
	}
}
