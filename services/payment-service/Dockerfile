FROM golang:1.26-alpine AS builder

WORKDIR /src/payment-service

COPY common /src/common
COPY payment-service/go.mod payment-service/go.sum ./
RUN go mod download

COPY payment-service .

RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/payment-service/payment-service .

EXPOSE 8085

CMD ["./payment-service"]
