#! /bin/bash

echo "Starting dev server..."
export DEV=true
export JRNL_ENABLE_LOGS=true
watchexec -r -e go -- go run . serve
