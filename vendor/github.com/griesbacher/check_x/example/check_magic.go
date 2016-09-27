package main

import (
	"github.com/griesbacher/check_x"
	"time"
)

func main() {
	//Set Plugin timeout, it will kill you and end properly
	check_x.StartTimeout(5 * time.Second)

	//Create a warning threshold
	warn, err := check_x.NewThreshold("10:")
	if err != nil {
		panic(err)
	}

	//Create a critical threshold
	crit, err := check_x.NewThreshold("@20:30")
	if err != nil {
		panic(err)
	}

	//do your magic
	measuredValue1 := 25.0
	measuredValue2 := 5.0

	//evaluate your magic
	status1 := check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(measuredValue1)
	status2 := check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(measuredValue2)

	//see what went wrong
	worstState, err := check_x.States{status1, status2}.GetWorst()
	if err != nil {
		panic(err)
	}

	//set some performancedata
	check_x.NewPerformanceData("foo", measuredValue1).Unit("B").Warn(warn).Crit(crit).Min(0).Max(100)
	check_x.NewPerformanceData("bar", measuredValue2).Unit("s").Min(0)

	//bring it to an end
	check_x.Exit(*worstState, "Made by check_x")
}
