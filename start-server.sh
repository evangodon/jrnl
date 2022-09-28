#! /bin/sh

PORT=30305

docker build -t jrnl .
docker run -it -p $PORT:8080  -v $HOME/.local/share/jrnl:/root/.local/share/jrnl jrnl
