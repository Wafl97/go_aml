package fsm

import (
	"strconv"
	"strings"

	"github.com/Wafl97/go_aml/util/logger"
	"github.com/Wafl97/go_aml/util/types"
)

var log logger.Logger

func FromString(str string) types.Option[FinitStateMachine] {
	log = logger.New("PARSER")
	log.Debug("Building model ...")

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
		log.Error("No initial state provided")
		return types.None[FinitStateMachine]()
	}
	log.Debug("Building complete")
	return types.Some(builder.Build())
}

func handleVariableDeclaration(line string, lineNumber int, builder *FsmBuilder) {
	varDef, isVarDef := strings.CutPrefix(line, "var")
	if !isVarDef {
		return
	}
	varName, varValue, isValidVarDef := strings.Cut(varDef, "=")
	if !isValidVarDef {
		log.Warnf("Bad variable declaration on line %d, invalid declaration ... skipping", lineNumber+1)
		return
	}
	varName = strings.TrimSpace(varName)
	varValue = strings.TrimSpace(varValue)
	if len(varName) == 0 {
		log.Warnf("Bad variable declaration on line %d, missing name ... skipping", lineNumber+1)
		return
	}
	if len(varValue) == 0 {
		log.Warnf("Bad variable declaration on line %d, missing value ... skipping", lineNumber+1)
		return
	}
	// try int first
	intValue, ierr := strconv.ParseInt(varValue, 10, 32)
	if ierr == nil {
		builder.variables.Set(varName, intValue)
		log.Debugf("Set Int")
		return
	}
	// then float
	floatValue, ferr := strconv.ParseFloat(varValue, 32)
	if ferr == nil {
		builder.variables.Set(varName, floatValue)
		log.Debugf("Set Float")
		return
	}
	// then bool
	boolValue, berr := strconv.ParseBool(varValue)
	if berr == nil {
		builder.variables.Set(varName, boolValue)
		log.Debugf("Set Bool")
		return
	}
	// if all else fails just have a string
	varValue = strings.Trim(varValue, "\"")
	builder.variables.Set(varName, varValue)
	log.Debugf("Set String")
}

func handleModelName(line string, lineNumber int, builder *FsmBuilder) {
	modelName, containsModelName := strings.CutPrefix(line, "model ")
	if containsModelName {
		if len(modelName) == 0 {
			log.Warnf("No model name given on line %d", lineNumber+1)
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
		log.Warnf("No name given to state on line %d ... skipping", lineNumber+1)
		return lineNumber
	}
	iterated := 0
	builder.Given(state, func(sb *StateBuilder) {
		for iterated = lineNumber + 1; iterated < len(*lines); iterated++ {
			line = (*lines)[iterated]
			if line == "}" { // end of state block
				break
			}
			checkLineIsTermination(line, iterated, sb, state)
			checkLineIsTransition(line, iterated, sb, state)
		}
		if len(sb.transitions) == 0 {
			log.Warnf("No valid transitions provided for state %s on line %d", state, lineNumber+1)
		}
	})
	if init == "init " {
		builder.Initial(state)
		log.Debugf("setting initial state %s", state)
	}
	return iterated
}

func checkLineIsTransition(line string, lineNumber int, sb *StateBuilder, state string) {
	event, nextState, isTransition := strings.Cut(line, "->")
	if !isTransition {
		return
	}
	event = strings.TrimSpace(event)
	nextState = strings.TrimSpace(nextState)
	if len(event) == 0 {
		log.Warnf("Bad transition on line %d, no event provided ... skipping", lineNumber+1)
		return
	}
	if len(nextState) == 0 {
		log.Warnf("Bad transition on line %d, no destination state provided ... skipping", lineNumber+1)
		return
	}
	sb.When(event, func(eb *EdgeBuilder) {
		eb.Then(nextState)
	})
	log.Debugf("%s on %s goto %s ... done", state, event, nextState)
}

func checkLineIsTermination(line string, lineNumber int, sb *StateBuilder, state string) {
	event, _, isTermination := strings.Cut(line, "-x")
	if !isTermination {
		return
	}
	event = strings.TrimSpace(event)
	if len(event) == 0 {
		log.Warnf("Bad termination on line %d, no event provided ... skipping", lineNumber+1)
		return
	}
	sb.When(event, func(eb *EdgeBuilder) { eb.End() })
	log.Debugf("%s on %s terminate ... done", state, event)
}
