# syntax = docker/dockerfile:1

FROM golang:1.18-alpine 

WORKDIR /backend-nobl9

COPY go.mod ./
RUN go mod download 

COPY *.go ./

RUN go build -o main .

EXPOSE 8000

CMD ["/backend-nobl9/main"]
