FROM golang:alpine

WORKDIR /app

RUN apk update
RUN apk add --virtual build-dependencies build-base gcc

#COPY go.mod /code/
COPY go.mod /app
RUN go mod download
RUN go mod verify


CMD ["go", "build"]
RUN go mod tidy
