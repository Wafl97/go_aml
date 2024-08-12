package fsm

import (
	"github.com/Wafl97/go_aml/fsm/mode"
	"github.com/Wafl97/go_aml/util/functions"
	"github.com/Wafl97/go_aml/util/types"
)

type EdgeMetaData struct {
	rawLine     string
	computation types.Option[string]
	condition   types.Option[string]
}

type Symbol uint

const (
	EQ Symbol = 0 // ==
	NE Symbol = 1 // !=
	GT Symbol = 2 // >
	GE Symbol = 3 // >=
	LT Symbol = 4 // <
	LE Symbol = 5 // <=
)

func (symbol Symbol) ToString() string {
	switch symbol {
	case EQ:
		return "=="
	case NE:
		return "!="
	case GT:
		return ">"
	case GE:
		return ">="
	case LT:
		return "<"
	case LE:
		return "<="
	default:
		return ""
	}
}

type Condition struct {
	left      string
	symbol    Symbol
	right     any
	valueType VariableType
}

type Computation struct {
	left      string
	operator  string
	right     any
	valueType VariableType
}

type Edge struct {
	terminate      mode.Mode
	resultingState types.Option[string]
	computation    types.Option[functions.Consumer[*Variables]]
	computation2   []Computation
	condition      types.Option[functions.Predicate[*Variables]]
	condition2     []Condition
	metaData       EdgeMetaData
}

func (edge *Edge) checkCondition(variables *Variables) (types.Option[string], mode.Mode) {
	next := edge.resultingState
	edge.condition.HasValue(func(p functions.Predicate[*Variables]) {
		if !p(variables) {
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
	edge.computation.HasValue(func(c functions.Consumer[*Variables]) {
		c(variables)
	})
}

type EdgeBuilder struct {
	terminate      mode.Mode
	resultingState types.Option[string]
	computation    types.Option[functions.Consumer[*Variables]] // DEPRECATED
	computation2   []Computation
	condition      types.Option[functions.Predicate[*Variables]] // DEPRECATED
	condition2     []Condition
	metaData       EdgeMetaData
}

func newEdgeBuilder() EdgeBuilder {
	return EdgeBuilder{
		terminate:      mode.CONTINUE,
		resultingState: types.None[string](),
		computation:    types.None[functions.Consumer[*Variables]](),
		computation2:   []Computation{},
		condition:      types.None[functions.Predicate[*Variables]](),
		condition2:     []Condition{},
		metaData: EdgeMetaData{
			computation: types.None[string](),
			condition:   types.None[string](),
		},
	}
}

func (builder *EdgeBuilder) MetaData(metaData string) *EdgeBuilder {
	builder.metaData.rawLine = metaData
	return builder
}

func (builder *EdgeBuilder) End() *EdgeBuilder {
	builder.terminate = mode.TERMINATE
	return builder
}

func (builder *EdgeBuilder) Then(state string) *EdgeBuilder {
	builder.resultingState = types.Some(state)
	return builder
}

// DEPRECATED
func (builder *EdgeBuilder) And(condition functions.Predicate[*Variables]) *EdgeBuilder {
	builder.condition = types.Some(condition)
	return builder
}

func (builder *EdgeBuilder) And2(conditions *[]Condition) *EdgeBuilder {
	builder.condition2 = append(builder.condition2, *conditions...)
	return builder
}

func (builder *EdgeBuilder) AndMeta(metaData string) *EdgeBuilder {
	builder.metaData.condition = types.Some(metaData)
	return builder
}

// DEPRECATED
func (builder *EdgeBuilder) Run(computation functions.Consumer[*Variables]) *EdgeBuilder {
	builder.computation = types.Some(computation)
	return builder
}

func (builder *EdgeBuilder) Run2(computations *[]Computation) *EdgeBuilder {
	builder.computation2 = append(builder.computation2, *computations...)
	return builder
}

func (builder *EdgeBuilder) RunMeta(metaData string) *EdgeBuilder {
	builder.metaData.computation = types.Some(metaData)
	return builder
}

func (builder *EdgeBuilder) build() Edge {
	return Edge{
		terminate:      builder.terminate,
		resultingState: builder.resultingState,
		computation:    builder.computation,
		computation2:   builder.computation2,
		condition:      builder.condition,
		condition2:     builder.condition2,
		metaData:       builder.metaData,
	}
}
