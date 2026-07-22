FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY services/common /src/common
COPY proto /proto
COPY services/notification-service/go.mod services/notification-service/go.sum ./
RUN go mod download
COPY services/notification-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/notification-service .
EXPOSE 8087
CMD ["./notification-service"]
