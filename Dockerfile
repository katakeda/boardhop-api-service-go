FROM golang:latest

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN make clean && make build

CMD ["./boardhop-api-service"]