package fsm

import (
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/logger"
)

type AutoEvent struct {
	conditions     Conditionals
	computations   Computational
	resultingState string
	//terminate      mode.Mode // DEPRECATED
	terminate2 bool
}

type State struct {
	logger              logger.Logger
	name                string
	defaultComputations Computational
	autoEvents          []AutoEvent
	transitions         map[string][]*Edge
	cache               map[string]any
}

func (state *State) GetTransitions() map[string][]*Edge {
	return state.transitions
}

func (state *State) fire(event string, variables *Variables) (*string, error) {
	arr, containsEvent := state.transitions[event]
	if !containsEvent {
		return nil, NewDeadlockErrorf("no transitions found for event %s", event)
	}
	state.logger.Debugf("Checking %d edge(s) ...", len(arr))
	for _, edge := range arr {
		res, err := edge.checkCondition(variables)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, NewDeadlockErrorf("fsm_state_fire: deadlock reached in state %s, no valid transitions", state.name)
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
	defaultComputations Computational
	autoEvents          []AutoEvent
	transitions         map[string][]*Edge
}

func newStateBuilder(state string) StateBuilder {
	return StateBuilder{
		logger: logger.New(state + "(Builder)"),
		name:   state,
		defaultComputations: Computational{
			FuncSignature: "func(event string)",
		},
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
	edgeBuilder.computation2.FuncSignature = "func()"
	edge := edgeBuilder.build()
	_, contains := builder.transitions[event]
	if contains {
		builder.transitions[event] = append(builder.transitions[event], &edge)
	} else {
		builder.transitions[event] = []*Edge{&edge}
	}
	return builder
}

func (builder *StateBuilder) AutoRun(computations *Computational) *StateBuilder {
	builder.defaultComputations = *computations
	return builder
}

func (builder *StateBuilder) AutoRunEvent(autoRunEvent AutoEvent) *StateBuilder {
	builder.autoEvents = append(builder.autoEvents, autoRunEvent)
	return builder
}
