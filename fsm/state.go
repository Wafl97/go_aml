package fsm

import (
	"github.com/Wafl97/go_aml/fsm/mode"
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/logger"
	"github.com/Wafl97/go_aml/util/types"
)

type AutoEvent struct {
	conditions     []Condition
	compuatations  []Computation
	resultingState string
	terminate      mode.Mode
}

type State struct {
	logger              logger.Logger
	name                string
	defaultComputations []Computation
	autoEvents          []AutoEvent
	transitions         map[string][]*Edge
	cache               map[string]any
}

func (state *State) GetTransitions() map[string][]*Edge {
	return state.transitions
}

func (state *State) fire(event string, variables *Variables) (types.Option[string], mode.Mode) {
	arr, containsEvent := state.transitions[event]
	if !containsEvent {
		return types.None[string](), mode.DEADLOCK
	}
	state.logger.Debugf("Checking %d edge(s) ...", len(arr))
	for _, edge := range arr {
		res, newMode := edge.checkCondition(variables)
		if res.IsSome() {
			return res, mode.CONTINUE
		}
		if newMode == mode.TERMINATE {
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

type StateBuilder struct {
	logger              logger.Logger
	name                string
	defaultComputations []Computation
	autoEvents          []AutoEvent
	transitions         map[string][]*Edge
}

func newStateBuilder(state string) StateBuilder {
	return StateBuilder{
		logger:      logger.New(state + "(Builder)"),
		name:        state,
		transitions: map[string][]*Edge{},
	}
}

func (builder *StateBuilder) build() State {
	return State{
		logger:              logger.New(builder.name),
		name:                builder.name,
		defaultComputations: builder.defaultComputations,
		autoEvents:          builder.autoEvents,
		transitions:         builder.transitions,
		cache:               map[string]any{},
	}
}

func (builder *StateBuilder) When(event string, f functions.Consumer[*EdgeBuilder]) *StateBuilder {
	edgeBuilder := newEdgeBuilder()
	f(&edgeBuilder)
	edge := edgeBuilder.build()
	_, contains := builder.transitions[event]
	if contains {
		builder.transitions[event] = append(builder.transitions[event], &edge)
	} else {
		builder.transitions[event] = []*Edge{&edge}
	}
	return builder
}

func (builder *StateBuilder) AutoRun(computations *[]Computation) *StateBuilder {
	builder.defaultComputations = *computations
	return builder
}

func (builder *StateBuilder) AutoRunEvent(autoRunEvent AutoEvent) *StateBuilder {
	builder.autoEvents = append(builder.autoEvents, autoRunEvent)
	builder.logger.Debug("A.R.E.")
	return builder
}
