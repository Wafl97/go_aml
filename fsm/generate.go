package fsm

import (
	"fmt"
	"os"
	"path"

	"github.com/Wafl97/go_aml/fsm/mode"
	"github.com/Wafl97/go_aml/util/logger"
)

const GENERATOR_VERSION = "v0.0.2"

var glog logger.Logger

const modFile = `// Generated by AML %s
module srcgen

go 1.21.6
`

func Generate(model *FinitStateMachine) {
	glog := logger.New("GENERATOR")
	glog.Infof("Generating code ...")
	os.MkdirAll("srcgen", os.ModeDir)

	generateFile(path.Join("srcgen", "go.mod"), fmt.Sprintf(modFile, GENERATOR_VERSION))
	generateFile(path.Join("srcgen", model.GetModelName()+".go"), generateCode(model))

	glog.Info("Generation complete")
}

func generateFile(fileName, fileContent string) {
	file, err := os.Create(fileName)
	if err != nil {
		glog.Error(err.Error())
	}
	file.Write([]byte(fileContent))
}

const codeStructure string = `/* Generated by AML %s */
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type State int
type Transition struct {
	condition 	   func() bool
	resultingState State
	function 	   func()
}
type StateNode struct {
	name        	   string	
	defaultComputation func(event string)
	transitions        map[string][]Transition
}

var ( /* VARIABLES */
%s)

const ( /* STATES */
	TERMINATION_STATE State = -1
%s)

var STATES []StateNode = []StateNode{
%s}

var CURRENT_STATE StateNode = STATES[STATE_%s]
	
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("State = %%s\n", CURRENT_STATE.name)
		switch event, err := reader.ReadString('\n'); err {
		case nil:
			event = strings.TrimSpace(event)
			switch state, success := CURRENT_STATE.transitions[event]; success {
			case false:
				if CURRENT_STATE.defaultComputation != nil {
					CURRENT_STATE.defaultComputation(event)
				}
				continue
			case true:
				for _, transition := range state {
					if transition.condition != nil && !transition.condition() {
						if CURRENT_STATE.defaultComputation != nil {
							CURRENT_STATE.defaultComputation(event)
						}
						continue
					}
					switch transition.resultingState {
					case TERMINATION_STATE:
						fmt.Print("Terminating")
						os.Exit(0)
					default:
						CURRENT_STATE = STATES[transition.resultingState]
						if transition.function != nil {
							transition.function()
						}
					}
				}
			}
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}
}
`

func generateCode(model *FinitStateMachine) string {
	var variables string
	for varName, varValue := range model.variables.values {
		variables += fmt.Sprintf("\t%s = %v\n", varName, varValue)
	}
	var states string
	stateCount := 0
	var transitions string
	for stateName, state := range model.states {
		states += fmt.Sprintf("\tSTATE_%s State = %d\n", stateName, stateCount)
		var defaultComputation string
		if state.defaultComputations == nil {
			defaultComputation = "nil"
		} else {
			defaultComputation = generateCompuation("func(event string)", *state.defaultComputations)
		}
		transitions += fmt.Sprintf("\t{\"%s\", %s, map[string][]Transition{ /* STATE_%s */\n", stateName, defaultComputation, stateName)
		for event, edges := range state.GetTransitions() {
			transitions += fmt.Sprintf("\t\t\"%s\": {\n", event)
			for _, edge := range edges {
				var resultState string
				switch edge.terminate {
				case mode.TERMINATE:
					resultState = "TERMINATION_STATE"
				default:
					resultState = fmt.Sprintf("STATE_%s", edge.resultingState.Get())
				}
				transitions += fmt.Sprintf("\t\t\t{%s, %s, %s}, /* %s */\n", generateCondition(edge), resultState, generateCompuation("func()", edge.computation2), edge.metaData.rawLine)
			}
			transitions += "\t\t},\n"
		}
		transitions += "\t}},\n"
		stateCount++
	}
	initialState := model.currentState.Get().GetName()
	return fmt.Sprintf(codeStructure, GENERATOR_VERSION, variables, states, transitions, initialState)
}

func generateCompuation(funcSignature string, computations []Computation) string {
	var computationString string
	switch len(computations) {
	case 0:
		computationString = "nil"
	default:
		computationString += fmt.Sprintf("%s {", funcSignature)
		for _, computation := range computations {
			computationString += fmt.Sprintf(" %s %s %v;", computation.left, computation.operator, computation.right)
		}
		computationString += " }"
	}
	return computationString
}

func generateCondition(edge *Edge) string {
	var contitionalString string
	switch len(edge.condition2) {
	case 0:
		contitionalString = "nil"
	default:
		contitionalString += "func() bool { return "
		for index, condition := range edge.condition2 {
			if condition.valueType == BOOL { // TODO: please for the love of god refactor this vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv
				if condition.right == "true" { //																							|
					if index == len(edge.condition2)-1 { //																					|
						contitionalString += fmt.Sprintf("%s", condition.left) //														|
					} else { //																												|
						contitionalString += fmt.Sprintf("%s && ", condition.left) // 													|
					} //																													|
				} else { //																													|
					if index == len(edge.condition2)-1 { // 																				|
						contitionalString += fmt.Sprintf("!%s", condition.left) // 														|
					} else { //																												|
						contitionalString += fmt.Sprintf("!%s && ", condition.left) //													|
					} //																													|
				} //																														|
			} else { //																														|
				if index == len(edge.condition2)-1 { //																						|
					contitionalString += fmt.Sprintf("%s %s %v", condition.left, condition.symbol.ToString(), condition.right) //		|
				} else { //																													|
					contitionalString += fmt.Sprintf("%s %s %v && ", condition.left, condition.symbol.ToString(), condition.right) // 	|
				} //																														|
			} // ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

		}
		contitionalString += " }"
	}
	return contitionalString
}
