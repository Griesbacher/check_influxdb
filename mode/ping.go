package mode

import (
	"fmt"
	"github.com/griesbacher/check_x"
	"github.com/influxdata/influxdb/client/v2"
	"time"
)

//Ping will be called on ping mode
func Ping(address, username, password string, insecureSkipVerify bool, warn, crit string) (err error) {
	warning, err := check_x.NewThreshold(warn)
	if err != nil {
		return
	}

	critical, err := check_x.NewThreshold(crit)
	if err != nil {
		return
	}

	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address, Username: username, Password: password, InsecureSkipVerify: insecureSkipVerify})
	if err != nil {
		return
	}
	defer c.Close()

	duration, version, err := c.Ping(time.Duration(1) * time.Second)
	if err != nil {
		return
	}

	value := duration.Seconds() * 1000
	state, err := check_x.States{
		check_x.Evaluator{
			Warning:  warning,
			Critical: critical,
		}.Evaluate(value),
	}.GetWorst()
	if err != nil {
		return
	}

	check_x.NewPerformanceData("request_duration", value).Unit("ms").Warn(warning).Crit(critical).Min(0)

	check_x.Exit(*state, fmt.Sprintf("InfluxDB version: %s", version))
	return nil
}
