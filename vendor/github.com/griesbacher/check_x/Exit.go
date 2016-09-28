package check_x

import (
	"fmt"
	"os"
)

//ErrorExit quits with unknown and the error message
func ErrorExit(err error) {
	Exit(Unknown, err.Error())
}

//Exit returns with the given returncode and message and optional performancedata
func Exit(state State, msg string) {
	LongExit(state, msg, "")
}

//LongExit returns with the given returncode and message and optional performancedata and long message
func LongExit(state State, msg, longMsg string) {
	if perf := PrintPerformanceData(); perf == "" {
		fmt.Printf("%s - %s\n%s", state.name, msg, longMsg)
	} else {
		fmt.Printf("%s - %s|%s\n%s", state.name, msg, perf, longMsg)
	}
	os.Exit(state.code)
}
