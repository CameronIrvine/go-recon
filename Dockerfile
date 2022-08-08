FROM golang:1.18-alpine

WORKDIR /app

RUN apk update
RUN apk add bind-tools

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /recon
