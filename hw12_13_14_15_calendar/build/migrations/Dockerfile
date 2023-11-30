FROM golang:1.20 as build

ENV CODE_DIR /go/src/migrations

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./migrations ${CODE_DIR}

RUN CGO_ENABLED=0 go build -o migrate *.go
RUN chmod +x migrate

CMD ["./migrate", "postgres"]