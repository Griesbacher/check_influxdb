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
)

func startTimeout() {
	if timeout != 0 {
		check_x.StartTimeout(time.Duration(timeout) * time.Second)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "check_influxdb"
	app.Usage = "Checks different influxdb stats\n   Copyright (c) 2016 Philip Griesbacher"
	app.Version = "0.0.1"
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
					Usage: "The diverence of series/measurements between now and x minutes. If a livestatus address is given, the overall state will switch to Warning if a core restart happend and due to that the metric got into Critical.",
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
					Usage: "Checks the bytes read and written the last x minutes",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.ReadWrite(address, username, password, warning, critical, timerange)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						flagFilter,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: read,write Bps (only read: 10, only write: ,10)",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: read,write Bps (only read: 10, only write: ,10)",
							Destination: &critical,
						},
						cli.IntFlag{
							Name:        "m",
							Usage:       "amount of minutes to look back",
							Destination: &timerange,
							Value:       1,
						},
					},
				},{
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
						flagFilter,
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
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "t",
			Usage:       "Seconds till check returns unkown, 0 to disable",
			Value:       10,
			Destination: &timeout,
		},
	}

	if err := app.Run(os.Args); err != nil {
		check_x.ErrorExit(err)
	}
}
