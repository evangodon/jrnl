#! /bin/sh

go build -o ./build/jrnl . 
echo "Building jrnl binary..."
pm2 restart build/jrnl -- serve --port 30305
