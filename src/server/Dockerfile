FROM golang:1.24.2-alpine AS base
WORKDIR /build
COPY . .


RUN go mod download && go get ./... && go build -o bin/server
RUN apk add --no-cache bash

EXPOSE 50051
EXPOSE 50052
EXPOSE 50053

RUN chmod +x run.sh

CMD ["./run.sh"]

