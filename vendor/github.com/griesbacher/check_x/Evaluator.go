package check_x

//Evaluator stores the warning and critical thresholds
type Evaluator struct {
	Warning  *Threshold
	Critical *Threshold
}

//Evaluate returns the stored values
func (t Evaluator) Evaluate(value float64) (result State) {
	if t.Critical != nil && !t.Critical.IsValueOK(value) {
		result = Critical
	} else if t.Warning != nil && !t.Warning.IsValueOK(value) {
		result = Warning
	} else {
		result = OK
	}
	return
}
