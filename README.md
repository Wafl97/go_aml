# A Modelling Language



## FSM

Finite State Machine. As part of the `.aml` format it is possible to model and simulte state machines. This is enables by setting the files syntax on the first line to `fsm`. This tells the program to use the fsm parser.
Within this syntax it is possible to declare variables and states, these variables can be used as guard conditions on transitions between states.
States are decalred as blocks, each contain all transitions from the owner state to another state.

## Usage

Note: the syntax declaration is required to be on line 1, the rest can be in any order 

```txt
syntax fsm

// give your state machine a name
model MY_MODEL

// presidence is int -> float -> bool -> string
var i = 10
var f = 0.5
var b = true
var s = some string

// everything within {} is a transition for the state
state STATE_1 {
  //event   -> resulting state 
    EVENT_1 -> STATE_2
  //event   && guard   -> state   && update
    EVENT_2 && i == 10 -> STATE_3 && i += 10
}

// init marks where the state machine will begin from
init state STATE_2 {
    EVENT_2 -> STATE_2
}

// -x marks a termination of the state machine
state STATE_3 {
    EVENT_3 -x
}

```

## Missing features

1. Conditional guards for transitions
2. Code execution (variable update) during transitions
3. Configure outputs
    1. Generate code
    2. Run random iterations
        1. Save result to file

## Future featues

1. Multiple state machines
    1. Different files for each state machine
    2. Communication between state machines
    


## Working on File

```go
func main() {
  fileBytes, err := os.ReadFile("../parser.aml")
  if err != nil {
    return
  }
  fileContent := string(fileBytes)
  chars := strings.Split(fileContent, "")
  for _, char := range chars {
    fmt.Println(char)
    handleEvent(char)
  }
}
```