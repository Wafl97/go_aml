package fsm

import (
	"fmt"
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

func handleModelName(line string, lineNumber int, builder *FsmBuilder) {
	modelName, containsModelName := strings.CutPrefix(line, "model ")
	if containsModelName {
		if len(modelName) == 0 {
			log.Warn(fmt.Sprintf("No model name given on line %d", lineNumber+1))
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
		log.Warn(fmt.Sprintf("No name given to state on line %d ... skipping", lineNumber+1))
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
			log.Warn(fmt.Sprintf("No transitions provided for state %s on line %d", state, lineNumber+1))
		}
	})
	if init == "init " {
		builder.Initial(state)
		log.Debug(fmt.Sprintf("setting initial state %s", state))
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
		log.Warn(fmt.Sprintf("Bad transition on line %d, no event provided ... skipping", lineNumber+1))
		return
	}
	if len(nextState) == 0 {
		log.Warn(fmt.Sprintf("Bad transition on line %d, no destination state provided ... skipping", lineNumber+1))
		return
	}
	sb.When(event, func(eb *EdgeBuilder) {
		eb.Then(nextState)
	})
	log.Debug(fmt.Sprintf("%s on %s goto %s ... done", state, event, nextState))
}

func checkLineIsTermination(line string, lineNumber int, sb *StateBuilder, state string) {
	event, _, isTermination := strings.Cut(line, "-x")
	if !isTermination {
		return
	}
	event = strings.TrimSpace(event)
	if len(event) == 0 {
		log.Warn(fmt.Sprintf("Bad termination on line %d, no event provided ... skipping", lineNumber+1))
		return
	}
	sb.When(event, func(eb *EdgeBuilder) { eb.End() })
	log.Debug(fmt.Sprintf("%s on %s terminate ... done", state, event))
}
