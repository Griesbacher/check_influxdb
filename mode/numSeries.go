package mode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_influxdb/helper"
	"github.com/griesbacher/check_x"
	"github.com/influxdata/influxdb/client/v2"
	"strconv"
)

func NumSeries(address, username, password, warning, critical, filterRegex string) (err error) {
	thresholds, err := helper.ParseCommaThresholds(warning, critical)
	if err != nil {
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address, Username: username, Password: password})
	if err != nil {
		return
	}
	defer c.Close()

	if filterRegex == "" {
		filterRegex = "."
	}

	q := client.NewQuery(
		fmt.Sprintf(`SELECT last(*) FROM "database" WHERE "database" =~ /%s/ GROUP BY  "database"`, filterRegex),
		"_internal", "s")
	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	data := map[string][]float64{}
	for _, r := range response.Results {
		for _, s := range r.Series {
			data[s.Tags["database"]] = []float64{}
			for _, i := range []int{
				helper.SliceIndex(len(s.Columns), func(i int) bool { return s.Columns[i] == "last_numMeasurements" }),
				helper.SliceIndex(len(s.Columns), func(i int) bool { return s.Columns[i] == "last_numSeries" }),
			} {
				f, err := s.Values[0][i].(json.Number).Float64()
				if err != nil {
					return err
				}
				data[s.Tags["database"]] = append(data[s.Tags["database"]], f)
			}
		}
	}

	states := check_x.States{}
	var okPrint bytes.Buffer
	var warnPrint bytes.Buffer
	var critPrint bytes.Buffer
	for database, values := range data {
		for i, v := range values {
			var w *check_x.Threshold
			var c *check_x.Threshold
			if len((*thresholds)["warning"])-1 >= i {
				w = (*thresholds)["warning"][i]
			}
			if len((*thresholds)["critical"])-1 >= i {
				c = (*thresholds)["critical"][i]
			}
			var typ string
			if typ = "s"; i == 0 {
				typ = "m"
			}
			label := fmt.Sprintf("%s__%s", typ, database)
			state := check_x.Evaluator{Warning: w, Critical: c}.Evaluate(v)
			states = append(states, state)

			writeToBuffer := func(buf *bytes.Buffer) {
				buf.WriteString(fmt.Sprintf("%s:%s ", label, strconv.FormatFloat(v, 'f', -1, 64)))
			}
			switch state {
			case check_x.OK:
				writeToBuffer(&okPrint)
				break
			case check_x.Warning:
				writeToBuffer(&warnPrint)
				break
			case check_x.Critical:
				writeToBuffer(&critPrint)
				break
			}

			check_x.NewPerformanceData(label, v).Min(0).Warn(w).Crit(c)
		}
	}
	var toPrint bytes.Buffer
	if critPrint.Len() > 0 {
		toPrint.WriteString("Critical: ")
		toPrint.WriteString(critPrint.String())
	}
	if warnPrint.Len() > 0 {
		toPrint.WriteString("Warning: ")
		toPrint.WriteString(warnPrint.String())
	}
	if okPrint.Len() > 0 {
		toPrint.WriteString("OK: ")
		toPrint.WriteString(okPrint.String())
	}
	worst, err := states.GetWorst()
	if err != nil {
		return
	}
	check_x.Exit(*worst, toPrint.String())

	return
}
