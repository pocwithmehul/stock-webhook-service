FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY ./stock-webhook-service /app/stock-webhook-service
COPY ./common-lib /app/common-lib
WORKDIR /app/stock-webhook-service
RUN go mod tidy
RUN go build -o /bin/stock-webhook-service ./...

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/stock-webhook-service /bin/stock-webhook-service
COPY --from=builder /app/stock-webhook-service/config /app/stock-webhook-service/config
WORKDIR /app/stock-webhook-service
ENTRYPOINT ["/bin/stock-webhook-service"]
