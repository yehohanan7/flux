#!/bin/bash
docker run --name mongo -p 27017:27017 -d mongo
go test $(go list ./... | grep -v /vendor/)
docker stop mongo
docker rm -v mongo
