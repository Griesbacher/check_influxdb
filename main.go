package main

import (
	"github.com/griesbacher/check_influxdb/mode"
	"github.com/griesbacher/check_x"
	"github.com/urfave/cli"
	"os"
	"time"
)

var (
	path        string
	address     string
	username    string
	password    string
	timeout     int
	warning     string
	critical    string
	filterRegex string
	timerange   int
	livestatus  string
	database    string
	query       string
	alias       string
	unknown2ok  bool
)

func startTimeout() {
	if timeout != 0 {
		check_x.StartTimeout(time.Duration(timeout) * time.Second)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "check_influxdb"
	app.Usage = "Checks different influxdb stats\n   Copyright (c) 2016 Philip Griesbacher\n   https://github.com/Griesbacher/check_influxdb"
	app.Version = "0.0.5"
	flagAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "InfluxDB address: Protocol + IP + Port",
		Destination: &address,
		Value:       "http://localhost:8086",
	}
	flagUsername := cli.StringFlag{
		Name:        "username",
		Destination: &username,
		Value:       "root",
	}
	flagPassword := cli.StringFlag{
		Name:        "password",
		Destination: &password,
		Value:       "root",
	}
	flagFilter := cli.StringFlag{
		Name:        "filter",
		Usage:       "regex to filter databases",
		Destination: &filterRegex,
	}
	app.Commands = []cli.Command{
		{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "check mode",
			Subcommands: []cli.Command{
				{
					Name:  "ping",
					Usage: "Tests if the influxdb is alive",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Ping(address, username, password, warning, critical)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: request duration in ms",
							Destination: &warning,
							Value:       "200",
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: request duration in ms",
							Destination: &critical,
							Value:       "400",
						},
					},
				}, {
					Name:  "disk",
					Usage: "Checks the disk size per databases",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Disk(path, warning, critical, filterRegex)
					},
					Flags: []cli.Flag{
						flagFilter,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: size in bytes",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: size in bytes",
							Destination: &critical,
						},
						cli.StringFlag{
							Name:        "path",
							Usage:       "path to the influxdb folder",
							Destination: &path,
							Value:       "/var/lib/influxdb",
						},
					},
				}, {
					Name:  "num_series",
					Usage: "The numbers of series/measurements",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.NumSeries(address, username, password, warning, critical, filterRegex)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						flagFilter,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: measurements,series (only measurements: 10, only series: ,10)",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: measurements,series (only measurements: 10, only series: ,10)",
							Destination: &critical,
						},
					},
				}, {
					Name:  "div_series",
					Usage: "The diverence of series/measurements between now and x minutes. If a livestatus address is given, the overall state will switch to Warning if a core restart happened and due to that the metric got into Critical.",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.DivSeries(address, username, password, warning, critical, filterRegex, livestatus, timerange)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						flagFilter,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: measurements,series (only measurements: 10, only series: ,10)",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: measurements,series (only measurements: 10, only series: ,10)",
							Destination: &critical,
						},
						cli.IntFlag{
							Name:        "m",
							Usage:       "amount of minutes to look back",
							Destination: &timerange,
							Value:       60,
						},
						cli.StringFlag{
							Name:        "l",
							Usage:       `livestatus address ("unix:/var/lib/nagios/rw/" or "tcp:localhost:6557")`,
							Destination: &livestatus,
						},
					},
				}, {
					Name:  "read_write",
					Usage: "Checks the bytes/operations read and written the last x minutes",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.ReadWrite(address, username, password, warning, critical, timerange)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: read(Bps),read(Ops),write(Bps),read(Ops) (only read: 10,10 only write: ,,10,10)",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: read(Bps),read(Ops),write(Bps),read(Ops) (only read: 10,10 only write: ,,10,10)",
							Destination: &critical,
						},
						cli.IntFlag{
							Name:        "m",
							Usage:       "amount of minutes to look back",
							Destination: &timerange,
							Value:       10,
						},
					},
				}, {
					Name:  "memory",
					Usage: "RSS in Byte",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Memory(address, username, password, warning, critical)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning in B",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical in B",
							Destination: &critical,
						},
					},
				}, {
					Name:  "old_series",
					Usage: "Returns a list of series older then x hours. This check makes only sense when the databases hast the tags: hostname and service - it's build for nagflux",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.OldSeries(address, database, username, password, warning, critical, timerange)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						cli.StringFlag{
							Name:        "database",
							Usage:       "Database to use",
							Destination: &database,
							Value:       "nagflux",
						},
						cli.IntFlag{
							Name:        "d",
							Usage:       "amount of hours",
							Destination: &timerange,
							Value:       24,
						},
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning in hours",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical in hours",
							Destination: &critical,
						},
					},
				}, {
					Name:  "query",
					Usage: "You could check a certain value from the database, but your query has to return only ONE value. Like 'select last(value) from metrics'",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Query(address, database, username, password, warning, critical, query, alias, unknown2ok)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						cli.StringFlag{
							Name:        "database",
							Usage:       "Database to use",
							Destination: &database,
							Value:       "nagflux",
						},
						cli.StringFlag{
							Name:        "q",
							Usage:       "query to be executed",
							Destination: &query,
						},
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning value",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical value",
							Destination: &critical,
						},
						cli.StringFlag{
							Name:        "a",
							Usage:       "alias, will replace the query within the output, if set",
							Destination: &alias,
						},
						cli.BoolFlag{
							Name:        "unknown2ok",
							Usage:       "If this flag is set, a query which would return unknown will now return ok (Default: false)",
							Destination: &unknown2ok,
						},
					},
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "t",
			Usage:       "Seconds till check returns unknown, 0 to disable",
			Value:       10,
			Destination: &timeout,
		},
	}

	if err := app.Run(os.Args); err != nil {
		check_x.ErrorExit(err)
	}
}
