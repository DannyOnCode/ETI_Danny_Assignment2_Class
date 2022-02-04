# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY *.html ./

RUN go build -o /ClassFrontEnd

EXPOSE 8040

CMD [ "/ClassFrontEnd" ]