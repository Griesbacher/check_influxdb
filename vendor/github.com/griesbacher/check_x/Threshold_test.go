package check_x

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

var toString = map[Threshold]string{
	Threshold{lower: 1.2, upper: 3.4, outside: true}:    "1.2 : 3.4 , outside: true",
	Threshold{lower: 1, upper: 3.4, outside: true}:      "1 : 3.4 , outside: true",
	Threshold{lower: 1.2, upper: 3, outside: true}:      "1.2 : 3 , outside: true",
	Threshold{lower: 1.2, upper: 3.4, outside: false}:   "1.2 : 3.4 , outside: false",
	Threshold{lower: -1.2, upper: -3.4, outside: false}: "-1.2 : -3.4 , outside: false",
}

func TestThreshold_String(t *testing.T) {
	for th, s := range toString {
		if fmt.Sprint(th) != s {
			t.Errorf("Got: %s expected: %s", th, s)
		}
	}
}

var stringToThreshold = []struct {
	input     string
	threshold *Threshold
	err       error
}{
	{"-3.4", &Threshold{input: "-3.4", lower: 0, upper: -3.4, outside: true}, nil},
	{" 3.4", &Threshold{input: "3.4", lower: 0, upper: 3.4, outside: true}, nil},
	{"3", &Threshold{input: "3", lower: 0, upper: 3, outside: true}, nil},
	{"-3", &Threshold{input: "-3", lower: 0, upper: -3, outside: true}, nil},
	{"foo", nil, errors.New("")},
	{"3,4", nil, errors.New("")},

	{" -3.4:", &Threshold{input: "-3.4:", lower: -3.4, upper: math.MaxFloat64, outside: true}, nil},
	{"3.4:", &Threshold{input: "3.4:", lower: 3.4, upper: math.MaxFloat64, outside: true}, nil},
	{"-3:", &Threshold{input: "-3:", lower: -3, upper: math.MaxFloat64, outside: true}, nil},
	{"3:", &Threshold{input: "3:", lower: 3, upper: math.MaxFloat64, outside: true}, nil},
	{"3,1:", nil, errors.New("")},

	{"~:-3.4 ", &Threshold{input: "~:-3.4", lower: minFloat64, upper: -3.4, outside: true}, nil},
	{"~:3.4", &Threshold{input: "~:3.4", lower: minFloat64, upper: 3.4, outside: true}, nil},
	{"~:-3", &Threshold{input: "~:-3", lower: minFloat64, upper: -3, outside: true}, nil},
	{" ~:3", &Threshold{input: "~:3", lower: minFloat64, upper: 3, outside: true}, nil},
	{"~:3,1", nil, errors.New("")},

	{"-1.2:-3.4", nil, FirstBiggerThenSecondError},
	{"3:2", nil, FirstBiggerThenSecondError},
	{"1.2:3.4", &Threshold{input: "1.2:3.4", lower: 1.2, upper: 3.4, outside: true}, nil},
	{"-1.2:3.4", &Threshold{input: "-1.2:3.4", lower: -1.2, upper: 3.4, outside: true}, nil},
	{"-3.4:-1.2", &Threshold{input: "-3.4:-1.2", lower: -3.4, upper: -1.2, outside: true}, nil},
	{"1.2:3", &Threshold{input: "1.2:3", lower: 1.2, upper: 3, outside: true}, nil},
	{"1:3", &Threshold{input: "1:3", lower: 1, upper: 3, outside: true}, nil},
	{"1:3.4", &Threshold{input: "1:3.4", lower: 1, upper: 3.4, outside: true}, nil},
	{"1,2:3,4", nil, errors.New("")},

	{"@-1.2:-3.4", nil, FirstBiggerThenSecondError},
	{" @3:2", nil, FirstBiggerThenSecondError},
	{"@1.2:3.4", &Threshold{input: "@1.2:3.4", lower: 1.2, upper: 3.4, outside: false}, nil},
	{"@-1.2:3.4 ", &Threshold{input: "@-1.2:3.4", lower: -1.2, upper: 3.4, outside: false}, nil},
	{"@-3.4:-1.2", &Threshold{input: "@-3.4:-1.2", lower: -3.4, upper: -1.2, outside: false}, nil},
	{"@1.2:3", &Threshold{input: "@1.2:3", lower: 1.2, upper: 3, outside: false}, nil},
	{"@1:3", &Threshold{input: "@1:3", lower: 1, upper: 3, outside: false}, nil},
	{"@1:3.4", &Threshold{input: "@1:3.4", lower: 1, upper: 3.4, outside: false}, nil},
	{"@1,2:3,4", nil, errors.New("")},
}

func TestNewThreshold(t *testing.T) {
	for i, data := range stringToThreshold {
		t_got, err := NewThreshold(data.input)
		if (err == nil) != (data.err == nil) {
			t.Error(i, err, data.err)
		}
		if (t_got != nil && data.threshold != nil) && *t_got != *data.threshold {
			t.Errorf("%d - Got: %s expected: %s", i, *t_got, *data.threshold)
		}
	}
}

var thresholdBorders = []struct {
	threshold string
	value     float64
	expected  bool
}{
	{"10", -1, false},
	{"10", 0, true},
	{"10", 1, true},
	{"10", 10, true},
	{"10", 11, false},

	{"10:", -1, false},
	{"10:", 9, false},
	{"10:", 10, true},
	{"10:", 11, true},

	{"~:10", 11, false},
	{"~:10", 10, true},
	{"~:10", 9, true},
	{"~:10", -1, true},

	{"10:20", -1, false},
	{"10:20", 9, false},
	{"10:20", 10, true},
	{"10:20", 11, true},
	{"10:20", 19, true},
	{"10:20", 20, true},
	{"10:20", 21, false},

	{"@10:20", -1, true},
	{"@10:20", 9, true},
	{"@10:20", 10, false},
	{"@10:20", 11, false},
	{"@10:20", 19, false},
	{"@10:20", 20, false},
	{"@10:20", 21, true},
}

func TestThreshold_IsValueOK(t *testing.T) {
	for i, data := range thresholdBorders {
		th, err := NewThreshold(data.threshold)
		if err != nil {
			t.Error(i, "There should be no error ", err)
		}
		result := th.IsValueOK(data.value)
		if result != data.expected {
			t.Errorf("%d - Expected: %t, got: %t", i, result, data.expected)
		}
	}
}
