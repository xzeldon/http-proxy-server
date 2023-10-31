#!/bin/bash

rm -rf ./bin && go build -ldflags "-s -w" -o ./bin/http-proxy-server main.go