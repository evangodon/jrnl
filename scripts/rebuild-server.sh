#! /bin/sh

go build -o ./build/jrnl . 
pm2 restart build/jrnl -- serve --port 30305
