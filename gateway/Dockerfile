FROM golang:1.26-alpine AS builder
WORKDIR /src/app
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download
COPY gateway .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/app/gateway .
EXPOSE 8080
CMD ["./gateway"]
