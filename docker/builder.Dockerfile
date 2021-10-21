FROM golang:1.17.2

WORKDIR /dripapp

COPY . .

RUN go get ./...

RUN make build-go