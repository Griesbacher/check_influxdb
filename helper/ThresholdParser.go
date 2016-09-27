package helper

import (
	"github.com/griesbacher/check_x"
	"strings"
)

func ParseCommaThresholds(warning, critical string) (*map[string][]*check_x.Threshold, error) {
	thresholds := map[string][]*check_x.Threshold{
		"warning":  []*check_x.Threshold{},
		"critical": []*check_x.Threshold{},
	}
	for i, s := range []string{warning, critical} {
		var index string
		if index = "critical"; i == 0 {
			index = "warning"
		}
		for _, t := range strings.Split(s, ",") {
			warn, err := check_x.NewThreshold(t)
			if err != nil {
				return nil, err
			}
			thresholds[index] = append(thresholds[index], warn)
		}
	}
	return &thresholds, nil
}
