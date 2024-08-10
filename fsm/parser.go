package fsm

import (
	"strconv"
	"strings"

	"github.com/Wafl97/go_aml/util/logger"
	"github.com/Wafl97/go_aml/util/types"
)

var plog logger.Logger

func FromString(str string) types.Option[FinitStateMachine] {
	plog = logger.New("PARSER")
	plog.Info("Building model ...")

	lines := strings.Split(str, "\r\n")

	builder := NewFsmBuilder()
	for lineNumber := 0; lineNumber < len(lines); lineNumber++ {
		line := strings.TrimSpace(lines[lineNumber])
		// skip comment lines
		if strings.HasPrefix(line, "//") {
			continue
		}

		// only attempt to find model name if none is alreay found
		if len(builder.modelName) == 0 {
			handleModelName(line, lineNumber, &builder)
		}

		// cast order: int -> float -> bool -> string
		handleVariableDeclaration(line, lineNumber, &builder)

		// set the new line number from iterating over the state definition
		lineNumber = handleStateDef(line, lineNumber, &builder, &lines)
	}

	if builder.initialState.IsNone() {
		plog.Error("No initial state provided")
		return types.None[FinitStateMachine]()
	}
	plog.Info("Building complete")
	return types.Some(builder.Build())
}

func handleVariableDeclaration(line string, lineNumber int, builder *FsmBuilder) {
	varDef, isVarDef := strings.CutPrefix(line, "var")
	if !isVarDef {
		return
	}
	varName, varValue, isValidVarDef := strings.Cut(varDef, "=")
	if !isValidVarDef {
		plog.Warnf("Bad variable declaration on line %d, invalid declaration ... skipping", lineNumber+1)
		return
	}
	varName = strings.TrimSpace(varName)
	varValue = strings.TrimSpace(varValue)
	if len(varName) == 0 {
		plog.Warnf("Bad variable declaration on line %d, missing name ... skipping", lineNumber+1)
		return
	}
	if len(varValue) == 0 {
		plog.Warnf("Bad variable declaration on line %d, missing value ... skipping", lineNumber+1)
		return
	}
	// try int first
	intValue, ierr := strconv.ParseInt(varValue, 10, 32)
	if ierr == nil {
		builder.variables.Set(varName, intValue)
		builder.variables.SetType(varName, INT)
		plog.Debugf("Set Int")
		return
	}
	// then float
	floatValue, ferr := strconv.ParseFloat(varValue, 32)
	if ferr == nil {
		builder.variables.Set(varName, floatValue)
		builder.variables.SetType(varName, FLOAT)
		plog.Debugf("Set Float")
		return
	}
	// then bool
	boolValue, berr := strconv.ParseBool(varValue)
	if berr == nil {
		builder.variables.Set(varName, boolValue)
		builder.variables.SetType(varName, BOOL)
		plog.Debugf("Set Bool")
		return
	}
	// if all else fails just have a string
	//varValue = strings.Trim(varValue, "\"")
	builder.variables.Set(varName, varValue)
	builder.variables.SetType(varName, STRING)
	plog.Debugf("Set String")
}

func handleModelName(line string, lineNumber int, builder *FsmBuilder) {
	modelName, containsModelName := strings.CutPrefix(line, "model ")
	if containsModelName {
		if len(modelName) == 0 {
			plog.Warnf("No model name given on line %d", lineNumber+1)
			return
		}
		builder.modelName = modelName
	}
}

func handleStateDef(line string, lineNumber int, builder *FsmBuilder, lines *[]string) int {
	init, state, containsState := strings.Cut(line, "state ")
	if !containsState {
		return lineNumber
	}
	state = strings.TrimSpace(strings.Trim(state, "{"))
	if len(state) == 0 {
		plog.Warnf("No name given to state on line %d ... skipping", lineNumber+1)
		return lineNumber
	}
	iterated := 0
	builder.Given(state, func(sb *StateBuilder) {
		for iterated = lineNumber + 1; iterated < len(*lines); iterated++ {
			line = (*lines)[iterated]
			if line == "}" { // end of state block
				break
			}
			line = strings.TrimSpace(line)
			if checkLineIsTermination(line, iterated, sb, state) {
				continue
			}
			if checkLineIsTransition(line, iterated, sb, state, builder) {
				continue
			}
			plog.Warnf("Bad transition on line %d, invalid declaration ... skipping", iterated+1)
		}
		if len(sb.transitions) == 0 {
			plog.Warnf("No valid transitions provided for state %s on line %d", state, lineNumber+1)
		}
	})
	if init == "init " {
		builder.Initial(state)
		plog.Debugf("setting initial state %s", state)
	}
	return iterated
}

func checkLineIsTransition(line string, lineNumber int, sb *StateBuilder, state string, builder *FsmBuilder) bool {
	event, nextState, isTransition := strings.Cut(line, "->")
	if !isTransition {
		return false
	}
	event, conditionsString, isConditional := strings.Cut(event, "(")
	nextState, computationString, hasComputation := strings.Cut(nextState, "(")
	event = cleanEventString(event)
	nextState = strings.TrimSpace(nextState)
	if len(event) == 0 {
		plog.Warnf("Bad transition on line %d, no event provided ... skipping", lineNumber+1)
		return true
	}
	if len(nextState) == 0 {
		plog.Warnf("Bad transition on line %d, no destination state provided ... skipping", lineNumber+1)
		return true
	}
	sb.When(event, func(eb *EdgeBuilder) {
		eb.Then(nextState).MetaData(line)
		if isConditional {
			eb.And2(*parseCondition(strings.Trim(strings.TrimSpace(conditionsString), ")"), lineNumber, builder))
			eb.AndMeta(conditionsString)
		}
		if hasComputation {
			eb.Run2(*parseComputation(strings.Trim(strings.TrimSpace(computationString), ")"), lineNumber, builder))
			eb.RunMeta(computationString)
		}
	})
	plog.Debugf("%s on %s goto %s ... done", state, event, nextState)
	return true
}

func parseComputation(computationString string, lineNumber int, builder *FsmBuilder) *[]Computation {
	subComputations := strings.Split(computationString, ",")
	computations := make([]Computation, 0, len(subComputations))
	for _, subComputation := range subComputations {
		subComputation = strings.TrimSpace(subComputation)
		var computation Computation
		tokens := strings.SplitN(subComputation, " ", 3)
		if len(tokens) != 3 {
			plog.Warnf("Bad computation for transition on line %d, cannot infer computation (%s) ... skipping", lineNumber+1, subComputation)
			continue
		}
		computation.left = tokens[0]
		if _, isDeclared := builder.variables.values[computation.left]; !isDeclared {
			plog.Warnf("Bad computation for transition on line %d, variable '%s' is not declared ... skipping", lineNumber+1, computation.left)
			continue
		}
		computation.valueType = builder.variables.types[computation.left]
		computation.right = tokens[2]
		computations = append(computations, computation)
	}
	return &computations
}

func parseCondition(conditionString string, lineNumber int, builder *FsmBuilder) *[]Condition {
	subConditions := strings.Split(conditionString, ",")
	conditions := make([]Condition, 0, len(subConditions))
	for _, subCondition := range subConditions {
		subCondition = strings.TrimSpace(subCondition)
		var condition Condition
		tokens := strings.SplitN(subCondition, " ", 3)
		if len(tokens) != 3 {
			plog.Warnf("Bad condition for transition on line %d, cannot infer condition (%s) ... skipping", lineNumber+1, subCondition)
			continue
		}
		condition.left = tokens[0]
		if _, isDeclared := builder.variables.values[condition.left]; !isDeclared {
			plog.Warnf("Bad condition for transition on line %d, variable '%s' is not declared ... skipping", lineNumber+1, condition.left)
			continue
		}
		condition.right = tokens[2]
		condition.valueType = builder.variables.types[condition.left]
		switch tokens[1] {
		case "==":
			condition.symbol = EQ
		case "!=":
			condition.symbol = NE
		case ">=":
			condition.symbol = GE
		case "<":
			condition.symbol = GT
		case "<=":
			condition.symbol = LE
		case ">":
			condition.symbol = LT
		default:
			plog.Warnf("Bad condition in transition on line %d, invalid symbol ... skipping", lineNumber+1)
		}
		conditions = append(conditions, condition)
	}
	return &conditions
}

func checkLineIsTermination(line string, lineNumber int, sb *StateBuilder, state string) bool {
	event, isTermination := strings.CutSuffix(line, "-x")
	if !isTermination {
		return false
	}
	event = cleanEventString(event)
	if len(event) == 0 {
		plog.Warnf("Bad transition on line %d, no event provided ... skipping", lineNumber+1)
		return true
	}
	sb.When(event, func(eb *EdgeBuilder) {
		eb.End().MetaData(line)
	})
	plog.Debugf("%s on %s terminate ... done", state, event)
	return true
}

func cleanEventString(event string) string {
	return strings.Trim(strings.TrimSpace(event), "\"")
}
