#!/bin/bash
# Script to build and create release package
# https://github.com/smurfpandey/what-game

export GOOS=windows
go fmt && go build