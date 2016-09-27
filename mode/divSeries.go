package mode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/griesbacher/check_influxdb/helper"
	"github.com/griesbacher/check_x"
	"github.com/influxdata/influxdb/client/v2"
	"strconv"
	"time"
)

func DivSeries(address, username, password, warning, critical, filterRegex, livestatus string, timerange int) (err error) {
	thresholds, err := helper.ParseCommaThresholds(warning, critical)
	if err != nil {
		return
	}

	coreRestart := false
	if livestatus != "" {
		l, err := helper.NewLivestatus(livestatus)
		if err != nil {
			return err
		}
		liveResult, err := l.Query("GET status\nColumns: program_start\n\n")
		if err != nil {
			return err
		} else if len(*liveResult) == 0 {
			return errors.New("Livestatus did not return anything")
		}
		u, err := strconv.ParseInt((*liveResult)[0][0], 10, 64)
		if err != nil {
			return err
		}
		coreTime := time.Unix(u, 0)
		if int(time.Now().Sub(coreTime).Minutes()) < timerange {
			coreRestart = true
		}
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
		fmt.Sprintf(`SELECT last(*) FROM "database" WHERE "database" =~ /%s/ GROUP BY  "database";
SELECT last(*) FROM "database" WHERE "database" =~ /%s/ AND time < now() - %dm GROUP BY  "database"`, filterRegex, filterRegex, timerange),
		"_internal", "s")

	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	getSM := func(index int) (map[string][]float64, error) {
		data := map[string][]float64{}
		for _, s := range response.Results[index].Series {
			data[s.Tags["database"]] = []float64{}
			for _, i := range []int{2, 1} {
				f, err := s.Values[0][i].(json.Number).Float64()
				if err != nil {
					return data, err
				}
				data[s.Tags["database"]] = append(data[s.Tags["database"]], f)
			}
		}
		return data, nil
	}
	dataNew, err := getSM(0)
	if err != nil {
		return
	}
	dataOld, err := getSM(1)
	if err != nil {
		return
	}

	dataDiv := map[string][]float64{}
	for k, newValue := range dataNew {
		if oldValue, ok := dataOld[k]; ok {
			dataDiv[k] = []float64{}
			for i := 0; i < 2; i++ {
				dataDiv[k] = append(dataDiv[k], newValue[i]-oldValue[i])
			}
		} else {
			dataDiv[k] = newValue
		}
	}

	states := check_x.States{}
	var okPrint bytes.Buffer
	var warnPrint bytes.Buffer
	var critPrint bytes.Buffer
	for database, values := range dataDiv {
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

			check_x.NewPerformanceData(label, v).Warn(w).Crit(c)
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

	if *worst == check_x.Critical && coreRestart {
		check_x.Exit(check_x.Warning, toPrint.String())
	} else {
		check_x.Exit(*worst, toPrint.String())
	}

	return
}
