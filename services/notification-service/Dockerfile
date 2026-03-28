FROM golang:1.26-alpine AS builder

WORKDIR /src/notification-service

COPY common /src/common
COPY notification-service/go.mod notification-service/go.sum ./
RUN go mod download

COPY notification-service .

RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/notification-service/notification-service .

EXPOSE 8087

CMD ["./notification-service"]
