FROM golang:1.26-alpine AS builder

WORKDIR /src/order-service

COPY services/common /src/common
COPY proto /proto
COPY services/order-service/go.mod services/order-service/go.sum ./
RUN go mod download

COPY services/order-service .

RUN CGO_ENABLED=0 GOOS=linux go build -o order-service .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/order-service/order-service .

EXPOSE 8083 8084

CMD ["./order-service"]
