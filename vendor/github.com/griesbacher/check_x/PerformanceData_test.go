package check_x

import (
	"reflect"
	"testing"
)

var perfdataToString = []struct {
	f        func() PerformanceData
	expected string
}{
	{func() PerformanceData {
		return *NewPerformanceDataString("a", "1")
	}, "'a'=1;;;;"},
	{func() PerformanceData {
		warn, _ := NewThreshold("10:")
		return NewPerformanceDataString("a", "2").Warn(warn)
	}, "'a'=2;10:;;;"},
	{func() PerformanceData {
		warn, _ := NewThreshold("10:")
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceDataString("a", "3").Warn(warn).Crit(crit)
	}, "'a'=3;10:;@10:20;;"},
	{func() PerformanceData {
		warn, _ := NewThreshold("10:")
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceDataString("a", "3").Warn(warn).Crit(crit).Min(0)
	}, "'a'=3;10:;@10:20;0;"},
	{func() PerformanceData {
		warn, _ := NewThreshold("10:")
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceDataString("a", "4").Warn(warn).Crit(crit).Min(0).Max(100)
	}, "'a'=4;10:;@10:20;0;100"},
	{func() PerformanceData {
		warn, _ := NewThreshold("10:")
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceDataString("a", "5").Warn(warn).Crit(crit).Min(0).Max(100).Unit("C")
	}, "'a'=5C;10:;@10:20;0;100"},
	{func() PerformanceData {
		warn, _ := NewThreshold("10:")
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceData("a", 6).Warn(warn).Crit(crit).Min(0).Max(100).Unit("C")
	}, "'a'=6C;10:;@10:20;0;100"},
	{func() PerformanceData {
		warn, _ := NewThreshold("")
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceData("a", 6).Warn(warn).Crit(crit).Min(0).Max(100).Unit("C")
	}, "'a'=6C;;@10:20;0;100"},
	{func() PerformanceData {
		crit, _ := NewThreshold("@10:20")
		return NewPerformanceData("a", 6).Warn(nil).Crit(crit).Min(0).Max(100).Unit("C")
	}, "'a'=6C;;@10:20;0;100"},
}

func TestPerformanceData_toString(t *testing.T) {
	for i, data := range perfdataToString {
		result := data.f()
		resultString := result.toString()
		if resultString != data.expected {
			t.Errorf("%d - Expected: %s, got: %s", i, data.expected, resultString)
		}
		if !reflect.DeepEqual(p[i], result) {
			t.Errorf("%d - Expected: %s, got: %s", i, p[i], result)
		}
	}
}

func TestPrintPerformanceData(t *testing.T) {
	p = []PerformanceData{}
	NewPerformanceDataString("a", "1")
	NewPerformanceDataString("b", "2")
	expected := "'a'=1;;;;" + " " + "'b'=2;;;; "
	if expected != PrintPerformanceData() {
		t.Errorf("Expected: %s, got: %s", expected, PrintPerformanceData())
	}
}
