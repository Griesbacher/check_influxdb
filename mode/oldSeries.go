package mode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_x"
	"github.com/influxdata/influxdb/client/v2"
	"sort"
	"strconv"
	"time"
)

//OldSeries will check for old series
func OldSeries(address, database, username, password string, insecureSkipVerify bool, warning, critical string, timerange int) (err error) {
	warn, err := check_x.NewThreshold(warning)
	if err != nil {
		return
	}

	crit, err := check_x.NewThreshold(critical)
	if err != nil {
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address, Username: username, Password: password, InsecureSkipVerify: insecureSkipVerify})
	if err != nil {
		return
	}
	defer c.Close()

	q := client.NewQuery(`SELECT LAST(value) FROM metrics GROUP BY host,service`, database, "s")
	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	timeDiff := func(date time.Time) int {
		return int(time.Now().Sub(date).Hours())
	}

	data := map[string]map[string]time.Time{}
	for _, r := range response.Results {
		for _, s := range r.Series {
			unix, err := s.Values[0][0].(json.Number).Int64()
			if err != nil {
				return err
			}
			date := time.Unix(unix, 0)
			if timeDiff(date) > timerange {
				if _, ok := data[s.Tags["host"]]; !ok {
					data[s.Tags["host"]] = map[string]time.Time{}
				}
				data[s.Tags["host"]][s.Tags["service"]] = date
			}
		}
	}

	states := check_x.States{}
	var okPrint bytes.Buffer
	var warnPrint bytes.Buffer
	var critPrint bytes.Buffer

	var longOutput bytes.Buffer
	longOutput.WriteString("\n")

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	lastHost := ""
	for _, host := range keys {
		for service, date := range data[host] {
			diff := float64(timeDiff(date))
			state := check_x.Evaluator{Warning: warn, Critical: crit}.Evaluate(diff)
			states = append(states, state)

			writeToBuffer := func(buf *bytes.Buffer) {
				buf.WriteString(fmt.Sprintf("%s-%s-%sh ", host, service, strconv.FormatFloat(diff, 'f', -1, 64)))
			}
			if lastHost != host {
				longOutput.WriteString(host)
				longOutput.WriteString(":\n")
				lastHost = host
			}
			longOutput.WriteString("- ")
			longOutput.WriteString(service)
			longOutput.WriteString(": ")
			longOutput.WriteString(date.String())
			longOutput.WriteString("\n")

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

	check_x.LongExit(*worst, toPrint.String(), longOutput.String())

	return
}
