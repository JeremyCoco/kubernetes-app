# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as GO_Instance

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build ./cmd/web

EXPOSE 8181

CMD [ "./web" ]
