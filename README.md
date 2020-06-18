# Oratio
The main entry point of Milobella. We send it text and it returns vocalizable
answer.

## Prerequisites

- Having ``golang`` installed [instructions](https://golang.org/doc/install)

## Install

```bash
$ go build -o bin/oratio cmd/oratio/main.go
```

## Run
```bash
$ bin/oratio -c config/oratio.toml
```

## Example of request
```bash
$ curl -i -X POST http://localhost:9100/talk/text -d '{"text": "Quelle heure il est ? "}'
```
