FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY services/common /src/common
COPY proto /proto
COPY services/restaurant-service/go.mod services/restaurant-service/go.sum ./
RUN go mod download
COPY services/restaurant-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o restaurant-service .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/restaurant-service .
EXPOSE 8082 9082
CMD ["./restaurant-service"]
