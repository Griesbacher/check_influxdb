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

	app.Commands = []cli.Command{
		{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "check mode",
			Subcommands: []cli.Command{
				{
					Name:  "ping",
					Usage: "tests if the influxdb is alive",
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
					Usage: "checks the disk size per databases",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.Disk(path, warning, critical, filterRegex)
					},
					Flags: []cli.Flag{
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
						cli.StringFlag{
							Name:        "filter",
							Usage:       "regex to filter databases",
							Destination: &filterRegex,
						},
					},
				}, {
					Name:  "num_series",
					Usage: "the numbers of series/measurements",
					Action: func(c *cli.Context) error {
						startTimeout()
						return mode.NumSeries(address, username, password, warning, critical, filterRegex)
					},
					Flags: []cli.Flag{
						flagAddress,
						flagUsername,
						flagPassword,
						cli.StringFlag{
							Name:        "w",
							Usage:       "warning: series,measurements",
							Destination: &warning,
						},
						cli.StringFlag{
							Name:        "c",
							Usage:       "critical: series,measurements",
							Destination: &critical,
						},
						cli.StringFlag{
							Name:        "filter",
							Usage:       "regex to filter databases",
							Destination: &filterRegex,
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
