package check_x

import (
	"fmt"
	"os"
)

//ErrorExit quits with unknown and the error message
func ErrorExit(err error) {
	Exit(Unknown, err.Error())
}

//Exit returns with the given returncode and message and performancedata
func Exit(state State, msg string) {
	if perf := PrintPerformanceData(); perf == "" {
		fmt.Printf("%s - %s\n", state.name, msg)
	} else {
		fmt.Printf("%s - %s|%s\n", state.name, msg, perf)
	}
	os.Exit(state.code)
}
