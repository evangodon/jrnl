#! /bin/sh

git pull origin main
echo "Building jrnl binary..."
go build -o ./build/jrnl . 

cd app
echo "Install web depencencies..."
yarn install
echo "Building web bundle..."
yarn build


pm2 restart jrnl -- serve --port 30305
