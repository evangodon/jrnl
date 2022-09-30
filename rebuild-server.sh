#! /bin/sh

PORT=30305
NAME='jrnl-server'

docker stop "$NAME" 
docker rm "$NAME" 
docker build -t jrnl .
docker run -it -p $PORT:8080  -v $HOME/.local/share/jrnl:/root/.local/share/jrnl --name "$NAME"  jrnl
