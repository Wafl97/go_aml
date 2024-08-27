package test

import (
	"testing"

	"github.com/Wafl97/go_aml/fsm"
	"github.com/Wafl97/go_aml/runners"
	"github.com/Wafl97/go_aml/util/logger"
)

func TestVariables(t *testing.T) {
	vars := fsm.NewVariables()
	vars.Set("A", fsm.INT, 1)
	a, _ := vars.Get("A")

	if a != 1 {
		t.Log("Variables.Get Failed")
		t.Fail()
	}

}

const MODEL_NAME = "MODEL T"
const STATE_1 = "STATE 1"
const STATE_2 = "STATE 2"
const EVENT_1 = "EVENT 1"

func TestFSM(t *testing.T) {
	smb := fsm.NewFsmBuilder()
	smb.Name(MODEL_NAME).Given2(STATE_1, fsm.NewStateBuilder().When2(EVENT_1, fsm.NewEdgeBuilder().Then(STATE_2))).Initial(STATE_1)
	sm := smb.Build()
	if sm.GetModelName() != MODEL_NAME {
		t.Fail()
	}
	if sm.GetCurrentState() == nil || sm.GetCurrentState().GetName() != STATE_1 {
		t.Fail()
	}
	sm.Fire(EVENT_1)
}

func TestFullModel(t *testing.T) {
	logger.SetLogLevel("debug")
	log := logger.New("TESTING")
	tmb := fsm.NewFsmBuilder()
	tmb.
		Name("MODEL T").
		DeclareVar("INT", fsm.INT, 0).
		DeclareVar("FLOAT", fsm.FLOAT, 0.1).
		Given("S1", func(sb *fsm.StateBuilder) {
			sb.
				When2("TO-S2", fsm.NewEdgeBuilder().Then("S2")).
				When("TO-S3", func(eb *fsm.EdgeBuilder) {
					eb.Then("S3")
				}).
				When("TO-S4", func(eb *fsm.EdgeBuilder) {
					eb.Then("S4")
				}).
				When("TO-S5", func(eb *fsm.EdgeBuilder) {
					eb.Then("S5")
				})
		}).
		Given("S2", func(sb *fsm.StateBuilder) {
			sb.When("TO-S1", func(eb *fsm.EdgeBuilder) {
				eb.Then("S1")
			}).
				When("TO-S3", func(eb *fsm.EdgeBuilder) {
					eb.Then("S3")
				}).
				When("TO-S4", func(eb *fsm.EdgeBuilder) {
					eb.Then("S4")
				}).
				When("TO-S5", func(eb *fsm.EdgeBuilder) {
					eb.Then("S5")
				})
		}).
		Given("S3", func(sb *fsm.StateBuilder) {
			sb.
				When("TO-S1", func(eb *fsm.EdgeBuilder) {
					eb.Then("S1")
				}).
				When("TO-S2", func(eb *fsm.EdgeBuilder) {
					eb.Then("S2")
				}).
				When("TO-S4", func(eb *fsm.EdgeBuilder) {
					eb.Then("S4")
				}).
				When("TO-S5", func(eb *fsm.EdgeBuilder) {
					eb.Then("S5")
				})
		}).
		Given("S4", func(sb *fsm.StateBuilder) {
			sb.
				When("TO-S1", func(eb *fsm.EdgeBuilder) {
					eb.Then("S1")
				}).
				When("TO-S2", func(eb *fsm.EdgeBuilder) {
					eb.Then("S2")
				}).
				When("TO-S3", func(eb *fsm.EdgeBuilder) {
					eb.Then("S3")
				}).
				When("TO-S5", func(eb *fsm.EdgeBuilder) {
					eb.Then("S5")
				})
		}).
		Given("S5", func(sb *fsm.StateBuilder) {
			sb.
				When("TO-S1", func(eb *fsm.EdgeBuilder) {
					eb.Then("S1")
				}).
				When("TO-S2", func(eb *fsm.EdgeBuilder) {
					eb.Then("S2")
				}).
				When("TO-S3", func(eb *fsm.EdgeBuilder) {
					eb.Then("S3")
				}).
				When("TO-S4", func(eb *fsm.EdgeBuilder) {
					eb.Then("S4")
				}).
				When("TO-SX", func(eb *fsm.EdgeBuilder) {
					eb.Then("SX").Run2(tmb.NewComputations(func(variables fsm.Variables) {
						variables.Update("INT", fsm.ADD_ASSIGN, 1)
						variables.Update("FLOAT", fsm.ASSIGN, float64(10.0))
					}))
				})
		}).
		Given("SX", func(sb *fsm.StateBuilder) {
			sb.
				When("TO-S1", func(eb *fsm.EdgeBuilder) {
					eb.Then("S1")
				})
		}).
		Initial("S1")

	tm := tmb.Build()

	log.Infof("%v", tm)

	sum, err := runners.RunAsRandom(&tm, 1000)
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("Path: %v", sum.Path)
	log.Infof("Occurrences: %v", sum.Occurrences)
}

func FuzzVariables(f *testing.F) {
	vars := fsm.NewVariables()

	f.Add(float32(0.5), 1, false)
	f.Fuzz(func(t *testing.T, f1 float32, i1 int, b1 bool) {
		vars.Set("f", fsm.FLOAT, f1)
		vars.Set("i", fsm.INT, i1)
		vars.Set("b", fsm.BOOL, b1)

		if f2, _ := vars.Get("f"); f2 != f1 {
			t.Error("float failed")
		}
		if i2, _ := vars.Get("i"); i2 != i1 {
			t.Error("int failed")
		}
		if b2, _ := vars.Get("b"); b2 != b1 {
			t.Error("bool failed")
		}
	})
}
