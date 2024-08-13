package test

import (
	"testing"

	"github.com/Wafl97/go_aml/fsm"
)

func TestGenerateCondition(t *testing.T) {

	expectedConditionsStrings := []string{
		"func() bool { return a == 1 }",                // TEST 1
		"func() bool { return a >= 0.5 && b != \"\" }", // TEST 2
		"func() bool { return !a && b && c <= 5 }",     // TEST 3
	}

	conditions := [][]fsm.Condition{
		{ // TEST 1
			{
				Left:      "a",
				Symbol:    fsm.EQ,
				Right:     "1",
				ValueType: fsm.INT,
			},
		},
		{ // TEST 2
			{
				Left:      "a",
				Symbol:    fsm.GE,
				Right:     "0.5",
				ValueType: fsm.FLOAT,
			},
			{
				Left:      "b",
				Symbol:    fsm.NE,
				Right:     "\"\"",
				ValueType: fsm.STRING,
			},
		},
		{ // TEST 3
			{
				Left:      "a",
				Right:     "false",
				ValueType: fsm.BOOL,
			},
			{
				Left:      "b",
				Right:     "true",
				ValueType: fsm.BOOL,
			},
			{
				Left:      "c",
				Symbol:    fsm.LE,
				Right:     5,
				ValueType: fsm.INT,
			},
		},
	}

	actualConditionsStrings := make([]string, len(expectedConditionsStrings))

	for test := 0; test < len(actualConditionsStrings); test++ {
		actualConditionsStrings[test] = fsm.GenerateCondition(&conditions[test])
		if actualConditionsStrings[test] != expectedConditionsStrings[test] {
			t.Error("generateCondition: failed")
		}
	}
}
