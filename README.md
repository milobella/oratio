# Oratio
The main entry point of Milobella. We send it text and it returns vocalizable
answer.

## Prerequisites

- Having ``golang`` installed [instructions](https://golang.org/doc/install)
- Having ``go dep`` installed [instructins](https://golang.github.io/dep/docs/installation.html)

## Install

```
dep ensure
go build -o bin/oratio cmd/oratio/main.go
```

## Run
```
bin/oratio -c config/oratio.toml
```

## Example of request
```bash
curl -i -X POST http://localhost:9100/talk/text -d '{"text": "Quelle heure il est ? "}'
```

## CHANGELOGS
- [Application changelog](./CHANGELOG.md)
- [Helm chart changelog](./helm/oratio/CHANGELOG.md)