FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY services/common /src/common
COPY proto /proto
COPY services/payment-service/go.mod services/payment-service/go.sum ./
RUN go mod download
COPY services/payment-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/payment-service .
EXPOSE 8085
CMD ["./payment-service"]
