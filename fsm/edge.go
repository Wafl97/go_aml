package fsm

import (
	"github.com/Wafl97/go_aml/util/functions"
)

type EdgeMetaData struct {
	rawLine     string
	computation *string
	condition   *string
}

type Edge struct {
	terminate2     bool
	resultingState *string
	computation    functions.Consumer[*Variables] // DEPRECATED
	computation2   Computational
	//condition      functions.Predicate[*Variables] // DEPRECATED
	condition2 Conditionals
	metaData   EdgeMetaData
}

func (edge *Edge) checkCondition(variables *Variables) (*string, error) {
	if c := edge.condition2.function; c != nil {
		if !c(variables) {
			return nil, nil
		}
		edge.compute(variables)
		return edge.resultingState, nil
	}
	edge.compute(variables)
	if edge.terminate2 {
		return nil, ErrTerminated
	}
	return edge.resultingState, nil
}

func (edge *Edge) compute(variables *Variables) {
	if c := edge.computation; c != nil {
		c(variables)
	}
	if edge.computation2.Computations != nil {
		for _, computation := range edge.computation2.Computations {
			computation.Compute(variables)
		}
	}
}

type EdgeBuilder struct {
	terminate2     bool
	resultingState *string
	computation    functions.Consumer[*Variables] // DEPRECATED
	computation2   Computational
	condition      functions.Predicate[*Variables] // DEPRECATED
	condition2     []Condition
	metaData       EdgeMetaData
}

func NewEdgeBuilder() *EdgeBuilder {
	return &EdgeBuilder{
		terminate2:     false,
		resultingState: nil,
		computation:    nil,
		computation2: Computational{
			Computations: []Computation{},
		},
		condition:  nil,
		condition2: []Condition{},
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
	builder.condition = condition
	return builder
}

func (builder *EdgeBuilder) And2(conditions *[]Condition) *EdgeBuilder {
	builder.condition2 = *conditions
	return builder
}

func (builder *EdgeBuilder) AndMeta(metaData string) *EdgeBuilder {
	builder.metaData.condition = &metaData
	return builder
}

// DEPRECATED
func (builder *EdgeBuilder) Run(computation functions.Consumer[*Variables]) *EdgeBuilder {
	builder.computation = computation
	return builder
}

func (builder *EdgeBuilder) Run2(computations []Computation) *EdgeBuilder {
	builder.computation2.Computations = computations
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
		//condition:      builder.condition,
		condition2: NewConditionals(builder.condition2, builder.condition),
		metaData:   builder.metaData,
	}
}
