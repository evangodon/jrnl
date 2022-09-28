FROM golang:1.18-bullseye

ENV GO111MODULE=on

ENV APP_HOME /go/src/jrnl
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o jrnl .

CMD ["./jrnl", "serve"]
