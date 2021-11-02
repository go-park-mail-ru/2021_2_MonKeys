FROM golang:1.17.2

WORKDIR /dripapp

COPY . .

# RUN make linter && make build-go
RUN make build-go