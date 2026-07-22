FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY services/common /src/common
COPY proto /proto
COPY services/user-service/go.mod services/user-service/go.sum ./
RUN go mod download
COPY services/user-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/user-service .
EXPOSE 8088 9088
CMD ["./user-service"]
