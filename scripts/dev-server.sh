#! /bin/bash
#

OUT=tmp/bin/server
PORT=8090

export DEV=true
watchexec -r -e go -- "go build -o \"$OUT\" ./cmd/server/ && \"$OUT\"  --port \"$PORT\""
