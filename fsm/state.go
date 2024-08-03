package fsm

import (
	"fmt"
	"strconv"

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

func (state *State) fire(event string, variables *Variables) (types.Option[string], Termination) {
	arr, containsEvent := state.transitions[event]
	if !containsEvent {
		return types.None[string](), DONT
	}
	state.logger.Debug("Checking " + strconv.Itoa(len(arr)) + " edge(s) ...")
	for _, edge := range arr {
		res, terminate := edge.checkCondition(variables)
		if res.IsSome() {
			return res, terminate
		}
	}
	return types.None[string](), DONT
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
