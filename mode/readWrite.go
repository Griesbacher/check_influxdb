package mode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_influxdb/helper"
	"github.com/griesbacher/check_x"
	"github.com/griesbacher/check_x/Units"
	"github.com/influxdata/influxdb/client/v2"
)

//ReadWrite checks the bytes read and written the last x minutes
func ReadWrite(address, username, password string, insecureSkipVerify bool, warning, critical string, timerange int) (err error) {
	thresholds, err := helper.ParseCommaThresholds(warning, critical)
	if err != nil {
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address, Username: username, Password: password, InsecureSkipVerify: insecureSkipVerify})
	if err != nil {
		return
	}
	defer c.Close()

	timerangeInSeconds := timerange * 60
	q := client.NewQuery(
		fmt.Sprintf(`SELECT
( last("queryRespBytes") - first("queryRespBytes") ) / %d,
( last("queryReq") - first("queryReq") ) / %d,
( last("writeReqBytes") - first("writeReqBytes") ) / %d,
( last("writeReq") - first("writeReq") ) / %d
FROM "httpd" where time > now() - %dm;`,
			timerangeInSeconds, timerangeInSeconds, timerangeInSeconds, timerangeInSeconds, timerange),
		"_internal", "s")
	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	data := make([]float64, 4)
	for _, r := range response.Results {
		for _, s := range r.Series {
			for i := range data {
				d, err := s.Values[0][i+1].(json.Number).Float64()
				if err != nil {
					return err
				}
				data[i] = d
			}
		}
	}

	var toPrint bytes.Buffer
	states := check_x.States{}
	for i, v := range data {
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
		switch i {
		case 0:
			typ := "read_data"
			toPrint.WriteString(fmt.Sprintf("%s: %sps ", typ, Units.ByteSize(v)))
			check_x.NewPerformanceData(typ, v).Unit("Bps").Warn(w).Crit(c)
		case 1:
			typ := "read_ops"
			toPrint.WriteString(fmt.Sprintf("%s: %0.3fops ", typ, v))
			check_x.NewPerformanceData(typ, v).Unit("ops").Warn(w).Crit(c)
		case 2:
			typ := "write_data"
			toPrint.WriteString(fmt.Sprintf("%s: %sps ", typ, Units.ByteSize(v)))
			check_x.NewPerformanceData(typ, v).Unit("Bps").Warn(w).Crit(c)
		case 3:
			typ := "write_ops"
			toPrint.WriteString(fmt.Sprintf("%s: %0.3fops ", typ, v))
			check_x.NewPerformanceData(typ, v).Unit("ops").Warn(w).Crit(c)
		}
	}

	worst, err := states.GetWorst()
	if err != nil {
		return
	}
	check_x.Exit(*worst, toPrint.String())
	return
}
