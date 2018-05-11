package mode

import (
	"encoding/json"
	"fmt"
	"github.com/griesbacher/check_x"
	"github.com/influxdata/influxdb/client/v2"
	"strconv"
)

//Query will execute the given query and evaluate the result
func Query(address, database, username, password string, insecureSkipVerify bool, warning, critical, query, alias string, unknown2ok bool) (err error) {
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

	q := client.NewQuery(query, database, "s")
	response, err := c.Query(q)
	if err != nil {
		return
	} else if response.Error() != nil {
		return response.Error()
	}

	errorReturnCode := check_x.Unknown
	if unknown2ok {
		errorReturnCode = check_x.OK
	}

	if len(response.Results) != 1 {
		check_x.Exit(errorReturnCode, fmt.Sprintf("The amount of results is not 1 it is: %d", len(response.Results)))
	}
	if len(response.Results[0].Series) != 1 {
		check_x.Exit(errorReturnCode, fmt.Sprintf("The amount of series is not 1 it is: %d", len(response.Results[0].Series)))
	}
	if len(response.Results[0].Series[0].Values) != 1 {
		check_x.Exit(errorReturnCode, fmt.Sprintf("The amount of lines is not 1 it is: %d", len(response.Results[0].Series[0].Values)))
	}
	if len(response.Results[0].Series[0].Values[0]) != 2 {
		check_x.Exit(errorReturnCode, fmt.Sprintf("The amount of fields is not 2 it is: %d", len(response.Results[0].Series[0].Values[0])))
	}

	var result float64
	result, err = response.Results[0].Series[0].Values[0][1].(json.Number).Float64()

	check_x.NewPerformanceData("query", result).Warn(warn).Crit(crit)
	state := check_x.Evaluator{Warning: warn, Critical: warn}.Evaluate(result)

	resultAsString := strconv.FormatFloat(result, 'f', -1, 64)
	if alias == "" {
		check_x.Exit(state, fmt.Sprintf("Query: '%s' returned: '%s'", query, resultAsString))
	} else {
		check_x.Exit(state, fmt.Sprintf("Alias: '%s' returned: '%s'", alias, resultAsString))
	}
	return
}
