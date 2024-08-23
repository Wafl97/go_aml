package fsm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Wafl97/go_aml/util/logger"
)

var plog logger.Logger

func FromString(str string) (*FiniteStateMachine, error) {
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

		// only attempt to find model name if none is already found
		if len(builder.modelName) == 0 {
			handleModelName(line, lineNumber, &builder)
		}

		// cast order: int -> float -> bool -> string
		handleVariableDeclaration(line, lineNumber, &builder)

		// set the new line number from iterating over the state definition
		lineNumber = handleStateDef(line, lineNumber, &builder, &lines)
	}

	if builder.initialState == nil {
		plog.Error("No initial state provided")
		return nil, fmt.Errorf("fromString: invalid model")
	}
	plog.Info("Building complete")
	model := builder.Build()
	return &model, nil
}

func handleVariableDeclaration(line string, lineNumber int, builder *FiniteStateMachineBuilder) {
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
		builder.variables.Set(varName, INT, intValue)
		plog.Debugf("Set Int")
		return
	}
	// then float
	floatValue, ferr := strconv.ParseFloat(varValue, 32)
	if ferr == nil {
		builder.variables.Set(varName, FLOAT, floatValue)
		plog.Debugf("Set Float")
		return
	}
	// then bool
	boolValue, berr := strconv.ParseBool(varValue)
	if berr == nil {
		builder.variables.Set(varName, BOOL, boolValue)
		plog.Debugf("Set Bool")
		return
	}
	// if all else fails just have a string
	//varValue = strings.Trim(varValue, "\"")
	builder.variables.Set(varName, STRING, varValue)
	plog.Debugf("Set String")
}

func handleModelName(line string, lineNumber int, builder *FiniteStateMachineBuilder) {
	modelName, containsModelName := strings.CutPrefix(line, "model ")
	if containsModelName {
		if len(modelName) == 0 {
			plog.Warnf("No model name given on line %d", lineNumber+1)
			return
		}
		builder.modelName = modelName
	}
}

func handleStateDef(line string, lineNumber int, builder *FiniteStateMachineBuilder, lines *[]string) int {
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
			if len(line) == 0 {
				continue
			}
			if checkLineIsAutoRunTermination(line, lineNumber, sb, builder) {
				continue
			}
			if checkLineIsAutoComputation(line, lineNumber, sb, builder) {
				continue
			}
			if checkLineIsAutoRunEvent(line, lineNumber, sb, builder) {
				continue
			}
			if checkLineIsTermination(line, iterated, sb, state, builder) {
				continue
			}
			if checkLineIsTransition(line, iterated, sb, state, builder) {
				continue
			}
			plog.Warnf("Bad transition on line %d, invalid declaration ... skipping", iterated+1)
		}
		if len(sb.transitions) == 0 && len(sb.autoEvents) == 0 {
			plog.Warnf("No valid transitions or auto-events provided for state %s on line %d", state, lineNumber+1)
		}
	})
	if init == "init " {
		builder.Initial(state)
		plog.Debugf("setting initial state %s", state)
	}
	return iterated
}

func checkLineIsAutoComputation(line string, lineNumber int, sb *StateBuilder, builder *FiniteStateMachineBuilder) bool {
	autoComputation, isAutoComputation := strings.CutPrefix(line, ">>")
	if !isAutoComputation {
		return false
	}
	computations := parseComputation(autoComputation, lineNumber, builder)
	sb.AutoRun(computations)
	return true
}

func checkLineIsAutoRunTermination(line string, lineNumber int, sb *StateBuilder, builder *FiniteStateMachineBuilder) bool {
	autoEvent, isAutoEvent := strings.CutPrefix(line, "|>")
	if !isAutoEvent {
		return false
	}
	autoEvent, isTermination := strings.CutSuffix(autoEvent, "-x")
	if !isTermination {
		return false
	}
	var autoRunEvent AutoEvent
	autoRunEvent.conditions.conditions = *parseCondition(autoEvent, lineNumber, builder)
	autoRunEvent.terminate2 = true
	autoRunEvent.computations.FuncSignature = "func()"
	sb.AutoRunEvent(autoRunEvent)
	return true

}

func checkLineIsAutoRunEvent(line string, lineNumber int, sb *StateBuilder, builder *FiniteStateMachineBuilder) bool {
	autoEvent, isAutoEvent := strings.CutPrefix(line, "|>")
	if !isAutoEvent {
		return false
	}
	autoEvent, nextState, isValid := strings.Cut(autoEvent, "->")
	if !isValid {
		return false
	}
	nextState, computationString, hasComputation := strings.Cut(nextState, "(")

	nextState = strings.TrimSpace(nextState)
	if len(nextState) == 0 {
		plog.Warnf("Bad transition on line %d, no destination state provided ... skipping", lineNumber+1)
		return true
	}
	var autoRunEvent AutoEvent
	if hasComputation {
		autoRunEvent.computations.Computations = parseComputation(strings.Trim(strings.TrimSpace(computationString), ")"), lineNumber, builder)
	}
	autoRunEvent.conditions.conditions = *parseCondition(autoEvent, lineNumber, builder)
	autoRunEvent.terminate2 = false
	autoRunEvent.resultingState = nextState
	autoRunEvent.computations.FuncSignature = "func()"
	sb.AutoRunEvent(autoRunEvent)
	return true
}

func checkLineIsTermination(line string, lineNumber int, sb *StateBuilder, state string, builder *FiniteStateMachineBuilder) bool {
	event, isTermination := strings.CutSuffix(line, "-x")
	if !isTermination {
		return false
	}
	event, conditionsString, isConditional := strings.Cut(event, "(")
	event = cleanEventString(event)
	if len(event) == 0 {
		plog.Warnf("Bad transition on line %d, no event provided ... skipping", lineNumber+1)
		return true
	}
	sb.When(event, func(eb *EdgeBuilder) {
		eb.End().MetaData(line)
		if isConditional {
			eb.And2(parseCondition(strings.Trim(strings.TrimSpace(conditionsString), ")"), lineNumber, builder))
			eb.AndMeta(conditionsString)
		}
	})
	plog.Debugf("%s on %s terminate ... done", state, event)
	return true
}

func checkLineIsTransition(line string, lineNumber int, sb *StateBuilder, state string, builder *FiniteStateMachineBuilder) bool {
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
			eb.And2(parseCondition(strings.Trim(strings.TrimSpace(conditionsString), ")"), lineNumber, builder))
			eb.AndMeta(conditionsString)
		}
		if hasComputation {
			eb.Run2(parseComputation(strings.Trim(strings.TrimSpace(computationString), ")"), lineNumber, builder))
			eb.RunMeta(computationString)
		}
	})
	plog.Debugf("%s on %s goto %s ... done", state, event, nextState)
	return true
}

func parseComputation(computationString string, lineNumber int, builder *FiniteStateMachineBuilder) []Computation {
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
		computation.Left = tokens[0]
		if _, isDeclared := builder.variables.Get(computation.Left); !isDeclared {
			plog.Warnf("Bad computation for transition on line %d, variable '%s' is not declared ... skipping", lineNumber+1, computation.Left)
			continue
		}
		computation.ValueType = builder.variables.types[computation.Left]
		switch tokens[1] {
		case "=":
			computation.Operator = ASSIGN
		case "+=":
			computation.Operator = ADD_ASSIGN
		case "-=":
			computation.Operator = SUB_ASSIGN
		case "*=":
			computation.Operator = MUL_ASSIGN
		case "/=":
			computation.Operator = DIV_ASSIGN
		default:
			plog.Warnf("Bad computation in transition on line %d, invalid symbol (%s) ... skipping", lineNumber+1, tokens[1])
			continue
		}
		computation.Right = tokens[2]
		computations = append(computations, computation)
	}
	return computations
}

func parseCondition(conditionString string, lineNumber int, builder *FiniteStateMachineBuilder) *[]Condition {
	subConditions := strings.Split(conditionString, ",")
	conditions := make([]Condition, 0, len(subConditions))
	for _, subCondition := range subConditions {
		subCondition = strings.TrimSpace(subCondition)
		tokens := strings.SplitN(subCondition, " ", 3)
		if len(tokens) != 3 {
			plog.Warnf("Bad condition for transition on line %d, cannot infer condition (%s) ... skipping", lineNumber+1, subCondition)
			continue
		}
		if _, isDeclared := builder.variables.Get(tokens[0]); !isDeclared {
			plog.Warnf("Bad condition for transition on line %d, variable '%s' is not declared ... skipping", lineNumber+1, tokens[0])
			continue
		}
		var symbol LogicSymbol
		switch tokens[1] {
		case "==":
			symbol = EQUAL
		case "!=":
			symbol = NOT_EQUAL
		case ">=":
			symbol = GRATER_THAN_OR_EQUAL
		case ">":
			symbol = GRATER_THAN
		case "<=":
			symbol = LESS_THAN_OR_EQUAL
		case "<":
			symbol = LESS_THAN
		default:
			plog.Warnf("Bad condition in transition on line %d, invalid symbol (%s) ... skipping", lineNumber+1, tokens[1])
			continue
		}
		conditions = append(conditions, NewCondition(tokens[0], symbol, tokens[2], builder.variables.types[tokens[0]]))
	}
	return &conditions
}

func cleanEventString(event string) string {
	return strings.Trim(strings.TrimSpace(event), "\"")
}
