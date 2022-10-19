#! /bin/sh

git pull origin main
go build -o ./build/jrnl . 
echo "Building jrnl binary..."

cd app
yarn install
echo "Install web depencencies..."
yarn build
echo "Building web bundle..."


pm2 restart build/jrnl -- serve --port 30305
