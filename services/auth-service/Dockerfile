FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY services/common /src/common
COPY proto /proto
COPY services/auth-service/go.mod services/auth-service/go.sum ./
RUN go mod download
COPY services/auth-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/auth-service .
EXPOSE 8081
CMD ["./auth-service"]
