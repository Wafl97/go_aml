package fsm

import (
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/types"
)

type EdgeMetaData struct {
	rawLine     string
	computation *string
	condition   *string
}

type Edge struct {
	terminate2     bool
	resultingState *string
	computation    types.Option[functions.Consumer[*Variables]] // DEPRECATED
	computation2   Computational
	condition      types.Option[functions.Predicate[*Variables]] // DEPRECATED
	condition2     Conditionals
	metaData       EdgeMetaData
}

func (edge *Edge) checkCondition(variables *Variables) (*string, error) {
	next := edge.resultingState
	edge.condition.HasValue(func(p functions.Predicate[*Variables]) {
		if !p(variables) {
			next = nil
			return
		}
		edge.compute(variables)
	}).Else(func() {
		edge.compute(variables)
	})
	if edge.terminate2 {
		return nil, TerminateError
	}
	return next, nil
}

func (edge *Edge) compute(variables *Variables) {
	edge.computation.HasValue(func(c functions.Consumer[*Variables]) {
		c(variables)
	})
}

type EdgeBuilder struct {
	terminate2     bool
	resultingState *string
	computation    types.Option[functions.Consumer[*Variables]] // DEPRECATED
	computation2   Computational
	condition      types.Option[functions.Predicate[*Variables]] // DEPRECATED
	condition2     Conditionals
	metaData       EdgeMetaData
}

func newEdgeBuilder() EdgeBuilder {
	return EdgeBuilder{
		terminate2:     false,
		resultingState: nil,
		computation:    types.None[functions.Consumer[*Variables]](),
		computation2: Computational{
			Computations: []Computation{},
		},
		condition: types.None[functions.Predicate[*Variables]](),
		condition2: Conditionals{
			Conditions: []Condition{},
		},
		metaData: EdgeMetaData{
			computation: nil,
			condition:   nil,
		},
	}
}

func (builder *EdgeBuilder) MetaData(metaData string) *EdgeBuilder {
	builder.metaData.rawLine = metaData
	return builder
}

func (builder *EdgeBuilder) End() *EdgeBuilder {
	builder.terminate2 = true
	return builder
}

func (builder *EdgeBuilder) Then(state string) *EdgeBuilder {
	builder.resultingState = &state
	return builder
}

// DEPRECATED
func (builder *EdgeBuilder) And(condition functions.Predicate[*Variables]) *EdgeBuilder {
	builder.condition = types.Some(condition)
	return builder
}

func (builder *EdgeBuilder) And2(conditions *Conditionals) *EdgeBuilder {
	builder.condition2 = *conditions
	return builder
}

func (builder *EdgeBuilder) AndMeta(metaData string) *EdgeBuilder {
	builder.metaData.condition = &metaData
	return builder
}

// DEPRECATED
func (builder *EdgeBuilder) Run(computation functions.Consumer[*Variables]) *EdgeBuilder {
	builder.computation = types.Some(computation)
	return builder
}

func (builder *EdgeBuilder) Run2(computations *Computational) *EdgeBuilder {
	builder.computation2 = *computations
	return builder
}

func (builder *EdgeBuilder) RunMeta(metaData string) *EdgeBuilder {
	builder.metaData.computation = &metaData
	return builder
}

func (builder *EdgeBuilder) build() Edge {
	return Edge{
		terminate2:     builder.terminate2,
		resultingState: builder.resultingState,
		computation:    builder.computation,
		computation2:   builder.computation2,
		condition:      builder.condition,
		condition2:     builder.condition2,
		metaData:       builder.metaData,
	}
}
