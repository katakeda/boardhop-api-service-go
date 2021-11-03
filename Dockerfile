FROM golang:latest

COPY . /app
WORKDIR /app

RUN go install
RUN make clean && make build

CMD ["./boardhop-api-service"]