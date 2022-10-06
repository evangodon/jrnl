#! /bin/bash

echo "Starting dev server..."
export DEV=true
watchexec -r -e go -- go run . serve --port 8090
