FROM golang:latest

COPY . /app
WORKDIR /app

RUN go get -d -v ./...
RUN make clean && make build

CMD ["./boardhop-api-service"]