package check_x

import (
	"errors"
	"sort"
)

//State represents an Nagioskind returncode
type State struct {
	name string
	code int
}

//String prints the name of the state
func (s State) String() string {
	return s.name
}

var (
	//OK - returncode: 0
	OK = State{name: "OK", code: 0}
	//Warning - returncode: 1
	Warning = State{name: "Warning", code: 1}
	//Critical - returncode: 2
	Critical = State{name: "Critical", code: 2}
	//Unknown - returncode: 3
	Unknown = State{name: "Unknown", code: 3}
)

//States is a list of state
type States []State

var EmptyStatesError = errors.New("The given States do not contain an State")

//Len for Sort interface
func (s States) Len() int {
	return len(s)
}

//Less for Sort interface
func (s States) Less(i, j int) bool {
	return s[i].code < s[j].code
}

//Swap for Sort interface
func (s States) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s States) getSorted() error {
	if len(s) == 0 {
		return EmptyStatesError
	}
	if !sort.IsSorted(s) {
		sort.Sort(s)
	}
	return nil
}
func (s States) GetBest() (*State, error) {
	if err := s.getSorted(); err == nil {
		return &s[0], nil
	} else {
		return nil, err
	}
}
func (s States) GetWorst() (*State, error) {
	if err := s.getSorted(); err == nil {
		return &s[len(s)-1], nil
	} else {
		return nil, err
	}
}
