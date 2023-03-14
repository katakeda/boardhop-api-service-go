FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go install -buildvcs=false
RUN make clean && make build

FROM debian:bullseye-slim AS runtime
WORKDIR /app
RUN apt-get update -y \
    && apt-get install -y --no-install-recommends openssl ca-certificates \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*
COPY ./.secrets/google-credentials.json /app/google-credentials.json
COPY --from=builder /app/boardhop-api-service /app/boardhop-api-service
ENTRYPOINT ["./boardhop-api-service"]