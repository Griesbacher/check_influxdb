[![Circle CI](https://circleci.com/gh/Griesbacher/check_influxdb/tree/master.svg?style=svg)](https://circleci.com/gh/Griesbacher/check_influxdb/tree/master)
# check_influxdb
Monitoring Plugin to check the health of an InfluxDB

## Usage
### Global Options
```
$ go run main.go -h
NAME:
   check_influxdb - Checks different influxdb stats
   Copyright (c) 2016 Philip Griesbacher

USAGE:
   main.exe [global options] command [command options] [arguments...]

VERSION:
   0.0.3

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
     div_series  The diverence of series/measurements between now and x minutes. If a livestatus address is given, the overall state will switch to Warning if a core restart happend and due to that the metric got into Critical.
     read_write  Checks the bytes read and written the last x minutes
     memory      RSS in Byte
     old_series  Returns a list of series older then x hours. This check makes only sense when the databases hast the tags: hostname and service - it's build for nagflux

OPTIONS:
   --help, -h  show help
```