#!/bin/bash

set -e

GO111MODULE=off go get github.com/githubnemo/CompileDaemon
GO111MODULE=off go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

migrate -database=$DB_URL -path=migrations up
CompileDaemon --build="go build -o main cmd/api/main.go" --command=./main
