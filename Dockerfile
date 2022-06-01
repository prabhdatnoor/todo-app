FROM golang:alpine

WORKDIR /

#COPY go.mod /code/
COPY go.mod /

CMD ["go", "build"]
