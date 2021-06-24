# Gormungandr

## Introduction

Implementation of jormungandr in go as multiple micro services
Gormungandr is the little brother of [Jormungandr](https://github.com/CanalTP/navitia/tree/dev/source/jormungandr) inside the [navitia](https://github.com/CanalTP/navitia) project.<br>
It aims to be the front of navitia, like Jormungandr, but with better perfomance. Indeed, Gormungandr is implemented like a micro service in Golang.
This service is under construction and just **route_schedules** API is implemented.

## Architecture

### Written in

Gormungandr is written in **Go**.<br>

### APIs

The web API is powered by [gin](https://github.com/gin-gonic/gin).<br>
Supported API list :

- `/route_schedules` exposes the route_schedules API for Navitia
- `/status` exposes general information about the Gormungandr (This is not the status of Navitia)

###  TODO

To be compliant with Jormungandr, we have to

- Implement the disruptions handling

## Build

To build this project you need at least [go 1.17](https://golang.org/dl)<br>
Dependencies are handled by go modules as such it is recommended to not checkout this in your *GOPATH*.

To build the project you just need to run the following command, at the root of the project:

```shell
make build
```

If you want to run the tests:

``` shell
make test
```

## Run it

Gormungandr run with a bunch of input parameters:

```shell
./schedules -h
```

```
Usage of ./schedules:
      --auth-cache-timeout duration              timeout for cache on authentication calls to db
  -c, --connection-string string                 connection string to the jormungandr database (default "host=localhost user=navitia password=navitia dbname=jormungandr sslmode=disable")
      --json-log                                 enable json logging
      --kraken string                            zmq addr for kraken (default "tcp://localhost:3000")
      --listen string                            [IP]:PORT to listen (default ":8080")
      --log-level string                         log level: debug, info, warn, error (default "debug")
      --max-postresql-connections int            sets the maximum number of open connections to the database (default 20)
      --newrelic-appname string                  application name in new relic (default "gormungandr")
      --newrelic-license string                  license key new relic
      --pprof-listen string[="localhost:6060"]   address to listen for pprof. format: "IP:PORT"
  -r, --rabbitmq-dsn string                      connection uri for rabbitmq (default "amqp://guest:guest@localhost:5672/")
      --skip-auth                                disable authentication
      --skip-stats                               disable statistics
      --stats-exchange string                    exchange where to send stats (default "stat_persistor_exchange_topic")
      --timeout duration                         timeout for call to kraken (default 1s)
      --version                                  show version
```

Run Gormungandr and call `http://localhost:port/status`

Exemple:

```
# local run
./schedules --listen localhost:5000 --kraken ipc:///tmp/default_kraken --skip-auth true
```

## With Docker

Use the pre-built docker image: [navitia/schedules](https://hub.docker.com/r/navitia/schedules)
Several **tags** exists:

- release/latest - The main tag. This is the last stable version
- X.X.X - Each main version is tagged with a num
- master - The current branch of development

## Contribute

To contribute, create a Github PR from your fork, we will please to read your contribution.<br>
Don't forget to **lint** and **format** your code before to push, otherwise the CI will be merciless with you.<br>

Install the linter.

```shell
# Install linter
make linter-install
```

Now you can

```shell
# Run linter
make lint

# Run formatting
make fmt
```

