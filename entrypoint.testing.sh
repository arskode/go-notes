#!/bin/bash

set -e

GO111MODULE=off go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
migrate -database=postgresql://postgres:postgres@pg/postgres?sslmode=disable -path=migrations up
go test ./... -cover
yes | migrate -database=postgresql://postgres:postgres@pg/postgres?sslmode=disable -path=migrations down
