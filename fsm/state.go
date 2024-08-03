package fsm

import (
	"fmt"

	"github.com/Wafl97/go_aml/fsm/mode"
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/logger"
	"github.com/Wafl97/go_aml/util/types"
)

type State struct {
	logger      logger.Logger
	name        string
	transitions map[string][]*Edge
	cache       map[string]any
}

func (state *State) fire(event string, variables *Variables) (types.Option[string], mode.Mode) {
	arr, containsEvent := state.transitions[event]
	if !containsEvent {
		return types.None[string](), mode.DEADLOCK
	}
	state.logger.Debug(fmt.Sprintf("Checking %d edge(s) ...", len(arr)))
	for _, edge := range arr {
		res, newMode := edge.checkCondition(variables)
		if res.IsSome() {
			return res, mode.CONTINUE
		}
		if newMode == mode.TERMINATED {
			return res, newMode
		}
	}
	return types.None[string](), mode.DEADLOCK
}

func (state *State) GetEdgeTriggers() []string {
	cached, contains := state.cache["edge-triggers"]
	if contains {
		return cached.([]string)
	}
	cache := make([]string, 0, len(state.transitions))
	for k := range state.transitions {
		cache = append(cache, k)
	}
	state.cache["edge-triggers"] = cache
	return cache
}

func (state *State) GetName() string {
	return state.name
}

func (state State) ToString() string {
	return fmt.Sprintf("{ name: %s, transitions: %v }", state.name, state.transitions)
}

type StateBuilder struct {
	logger      logger.Logger
	name        string
	transitions map[string][]*Edge
}

func newStateBuilder(state string) StateBuilder {
	return StateBuilder{
		logger:      logger.New(state + "(Builder)"),
		name:        state,
		transitions: map[string][]*Edge{},
	}
}

func (sb *StateBuilder) build() State {
	return State{
		logger:      logger.New(sb.name),
		name:        sb.name,
		transitions: sb.transitions,
		cache:       map[string]any{},
	}
}

func (sb *StateBuilder) When(event string, f functions.Consumer[*EdgeBuilder]) *StateBuilder {
	edgeBuilder := newEdgeBuilder()
	f(&edgeBuilder)
	edge := edgeBuilder.build()
	_, contains := sb.transitions[event]
	if contains {
		sb.transitions[event] = append(sb.transitions[event], &edge)
	} else {
		sb.transitions[event] = []*Edge{&edge}
	}
	return sb
}
