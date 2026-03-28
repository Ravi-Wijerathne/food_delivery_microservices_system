FROM golang:1.26-alpine AS builder

WORKDIR /src/delivery-service

COPY common /src/common
COPY delivery-service/go.mod delivery-service/go.sum ./
RUN go mod download

COPY delivery-service .

RUN CGO_ENABLED=0 GOOS=linux go build -o delivery-service .

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /src/delivery-service/delivery-service .

EXPOSE 8086

CMD ["./delivery-service"]
