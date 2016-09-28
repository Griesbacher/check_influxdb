package mode

import (
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_x"
	"github.com/griesbacher/check_x/Units"
	"github.com/influxdata/influxdb/client/v2"
)

//Memory will check the RSS usage
func Memory(address, username, password, warning, critical string) (err error) {
	warn, err := check_x.NewThreshold(warning)
	if err != nil {
		return
	}

	crit, err := check_x.NewThreshold(critical)
	if err != nil {
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address, Username: username, Password: password})
	if err != nil {
		return
	}
	defer c.Close()

	q := client.NewQuery(`SELECT last("Sys") FROM "runtime"`, "_internal", "s")
	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	var mem float64
	for _, r := range response.Results {
		for _, s := range r.Series {
			mem, err = s.Values[0][1].(json.Number).Float64()
			if err != nil {
				return err
			}
		}
	}

	worst, err := check_x.States{
		check_x.Evaluator{
			Warning:  warn,
			Critical: crit,
		}.Evaluate(mem),
	}.GetWorst()
	if err != nil {
		return
	}

	check_x.NewPerformanceData("memory", mem).Unit("B").Warn(warn).Crit(crit)

	check_x.Exit(*worst, fmt.Sprintf("Memory usage: %s", Units.ByteSize(mem)))
	return
}
