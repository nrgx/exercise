#! /bin/bash

mkdir -p .cover
go test -coverprofile=.cover/coverage.out ./...; go tool cover -html=.cover/coverage.out
