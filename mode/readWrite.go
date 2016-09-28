package mode

import (
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_influxdb/helper"
	"github.com/griesbacher/check_x"
	"github.com/griesbacher/check_x/Units"
	"github.com/influxdata/influxdb/client/v2"
)

//ReadWrite checks the bytes read and written the last x minutes
func ReadWrite(address, username, password, warning, critical string, timerange int) (err error) {
	thresholds, err := helper.ParseCommaThresholds(warning, critical)
	if err != nil {
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address, Username: username, Password: password})
	if err != nil {
		return
	}
	defer c.Close()

	q := client.NewQuery(
		fmt.Sprintf(`SELECT last("queryRespBytes") / ( %d * 60 ) - first("queryRespBytes") / ( %d * 60 ),
last("writeReqBytes") / ( %d * 60 ) - first("writeReqBytes") / ( %d * 60 ) FROM "httpd" where time > now() - %dm;`,
			timerange, timerange, timerange, timerange, timerange),
		"_internal", "s")
	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	var write float64
	var read float64
	for _, r := range response.Results {
		for _, s := range r.Series {
			read, err = s.Values[0][1].(json.Number).Float64()
			if err != nil {
				return err
			}
			write, err = s.Values[0][2].(json.Number).Float64()
			if err != nil {
				return err
			}
		}
	}

	states := check_x.States{}
	for i, v := range []float64{read, write} {
		var w *check_x.Threshold
		var c *check_x.Threshold
		if len((*thresholds)["warning"])-1 >= i {
			w = (*thresholds)["warning"][i]
		}
		if len((*thresholds)["critical"])-1 >= i {
			c = (*thresholds)["critical"][i]
		}
		state := check_x.Evaluator{Warning: w, Critical: c}.Evaluate(v)
		states = append(states, state)
		var typ string
		if typ = "write"; i == 0 {
			typ = "read"
		}
		check_x.NewPerformanceData(typ, v).Unit("Bps").Warn(w).Crit(c)
	}

	worst, err := states.GetWorst()
	if err != nil {
		return
	}
	check_x.Exit(*worst, fmt.Sprintf("Read: %sps Write: %sps", Units.ByteSize(read), Units.ByteSize(write)))
	return
}
