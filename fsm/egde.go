package fsm

import (
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/types"
)

type Termination bool

const (
	DO   Termination = true
	DONT Termination = false
)

type Edge struct {
	terminate      Termination
	resultingState types.Option[string]
	computation    types.Option[functions.Consumer[Variables]]
	condition      types.Option[functions.Predicate[Variables]]
}

func NewEdge() Edge {
	return Edge{
		resultingState: types.None[string](),
		computation:    types.None[functions.Consumer[Variables]](),
		condition:      types.None[functions.Predicate[Variables]](),
	}
}

func (edge *Edge) checkCondition(variables *Variables) (types.Option[string], Termination) {
	next := edge.resultingState
	edge.condition.HasValue(func(p functions.Predicate[Variables]) {
		if !p(*variables) {
			next = types.None[string]()
			return
		}
		edge.compute(variables)
	}).Else(func() {
		edge.compute(variables)
	})
	return next, edge.terminate
}

func (edge *Edge) compute(variables *Variables) {
	edge.computation.HasValue(func(c functions.Consumer[Variables]) {
		c(*variables)
	})
}

type EdgeBuilder struct {
	terminate      Termination
	resultingState types.Option[string]
	computation    types.Option[functions.Consumer[Variables]]
	condition      types.Option[functions.Predicate[Variables]]
}

func newEdgeBuilder() EdgeBuilder {
	return EdgeBuilder{
		terminate:      DONT,
		resultingState: types.None[string](),
		computation:    types.None[functions.Consumer[Variables]](),
		condition:      types.None[functions.Predicate[Variables]](),
	}
}

func (edge *EdgeBuilder) End() *EdgeBuilder {
	edge.terminate = DO
	return edge
}

func (edge *EdgeBuilder) Then(state string) *EdgeBuilder {
	edge.resultingState = types.Some(state)
	return edge
}

func (edge *EdgeBuilder) And(condition functions.Predicate[Variables]) *EdgeBuilder {
	edge.condition = types.Some(condition)
	return edge
}

func (edge *EdgeBuilder) Run(computation functions.Consumer[Variables]) *EdgeBuilder {
	edge.computation = types.Some(computation)
	return edge
}

func (edge *EdgeBuilder) build() Edge {
	return Edge{
		terminate:      edge.terminate,
		resultingState: edge.resultingState,
		computation:    edge.computation,
		condition:      edge.condition,
	}
}
