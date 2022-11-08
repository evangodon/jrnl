#! /bin/sh

set -o errexit

git pull origin main
echo "Building jrnl binary..."
go build -o ./build/server ./cmd/server/

cd app
echo "Install web depencencies..."
yarn install
echo "Building web bundle..."
yarn build

pm2 restart jrnl -- serve --port 30305
