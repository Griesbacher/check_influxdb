[![Circle CI](https://circleci.com/gh/Griesbacher/check_influxdb/tree/master.svg?style=svg)](https://circleci.com/gh/Griesbacher/check_influxdb/tree/master)
# check_influxdb
Monitoring Plugin to check the health of an InfluxDB and data

## Usage
### Global Options
```
$ go run main.go -h
NAME:
   check_influxdb - Checks different influxdb stats
   Copyright (c) 2016 Philip Griesbacher
   https://github.com/Griesbacher/check_influxdb

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.5

COMMANDS:
     mode, m  check mode
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -t value       Seconds till check returns unkown, 0 to disable (default: 10)
   --help, -h     show help
   --version, -v  print the version
```

### Command options

```
$ go run main.go mode -h
NAME:
   check_influxdb mode - check mode

USAGE:
   check_influxdb mode command [command options] [arguments...]

COMMANDS:
     ping        Tests if the influxdb is alive
     disk        Checks the disk size per databases
     num_series  The numbers of series/measurements
     div_series  The diverence of series/measurements between now and x minutes. If a livestatus address is given, the overall state will switch to Warning if a core restart happened and due to that the metric got into Critical.
     read_write  Checks the bytes/operations read and written the last x minutes
     memory      RSS in Byte
     old_series  Returns a list of series older then x hours. This check makes only sense when the databases hast the tags: hostname and service - it's build for nagflux
     query       You could check a certain value from the database, but your query has to return only ONE value. Like 'select last(value) from metrics'
     
OPTIONS:
   --help, -h  show help
```

### Subcommand options

```
$ check_influxdb mode ping -h
NAME:
   check_influxdb mode ping - Tests if the influxdb is alive

USAGE:
   check_influxdb mode ping [command options] [arguments...]

OPTIONS:
   --address value   InfluxDB address: Protocol + IP + Port (default: "http://localhost:8086")
   --username value  (default: "root")
   --password value  (default: "root")
   -w value          warning: request duration in ms (default: "200")
   -c value          critical: request duration in ms (default: "400")

```