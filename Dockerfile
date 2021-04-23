FROM golang:latest

WORKDIR /app

COPY go.mod . 

COPY go.sum . 