#! /bin/sh

PORT=30305
IMAGE_NAME='jrnl'
CONTAINER_NAME='jrnl-server'

docker stop "$CONTAINER_NAME" 
docker rm "$CONTAINER_NAME" 
docker build -t "$IMAGE_NAME" .
docker run -p $PORT:8080  -v $HOME/.local/share/jrnl:/root/.local/share/jrnl --name "$NAME" --restart unless-stopped  "$IMAGE_NAME"
