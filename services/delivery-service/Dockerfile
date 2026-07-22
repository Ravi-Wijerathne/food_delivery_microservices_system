FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY services/common /src/common
COPY proto /proto
COPY services/delivery-service/go.mod services/delivery-service/go.sum ./
RUN go mod download
COPY services/delivery-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o delivery-service .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/delivery-service .
EXPOSE 8086
CMD ["./delivery-service"]
