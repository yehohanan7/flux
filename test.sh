#!/bin/bash
go test $(go list ./... | grep -v /vendor/) -logtostderr=true
