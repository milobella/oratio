![Build CI](https://github.com/milobella/oratio/workflows/Build%20CI/badge.svg)
![Deploy CI](https://github.com/milobella/oratio/workflows/Deploy%20CI/badge.svg)

# Oratio
The main entry point of Milobella. We send it text and it returns vocalizable answer.

## Prerequisites

- Having ``golang`` installed [instructions](https://golang.org/doc/install)

## Install

```bash
$ go build -o bin/oratio cmd/oratio/main.go
```

## Run
```bash
$ bin/oratio
```
> configuration ``config.toml`` will be checked in the following paths
> - /etc/oratio
> - $HOME/.oratio
> - .

A configuration example can be found in [config.toml](./config.toml).

## Examples of requests
### Talk to oratio
```bash
$ curl -i -X POST http://localhost:9100/talk/text -d '{"text": "Quelle heure il est ? "}'
```

### Register a new ability
```bash
$ curl -iv -X POST http://localhost:9100/abilities -d '{"name": "clock", "intents":["GET_TIME"], "host": "localhost", "port": 10300}'
```

### Get all registered abilities
```bash
$ curl -iv -X GET http://localhost:9100/abilities
```
