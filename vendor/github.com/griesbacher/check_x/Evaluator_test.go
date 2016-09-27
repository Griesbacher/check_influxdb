package check_x

import (
	"testing"
)

var thresholdsToState = []struct {
	warn  string
	crit  string
	value float64
	state State
}{
	{"10", "", 0, OK},
	{"10", "", 10, OK},
	{"10", "", -1, Warning},
	{"10", "", 11, Warning},
	{"", "10", 0, OK},
	{"", "10", 10, OK},
	{"", "10", -1, Critical},
	{"", "10", 11, Critical},
}

func TestEvaluator_Evaluate(t *testing.T) {
	for i, data := range thresholdsToState {
		var warn, crit *Threshold
		var err error
		if data.warn != "" {
			if warn, err = NewThreshold(data.warn); err != nil {
				t.Error(i, err)
			}
		}
		if data.crit != "" {
			if crit, err = NewThreshold(data.crit); err != nil {
				t.Error(i, err)
			}
		}

		if state := (Evaluator{Warning: warn, Critical: crit}.Evaluate(data.value)); state != data.state {
			t.Errorf("%d - Got state: %s - expected: %s", i, state, data.state)
		}
	}
}
