FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o modem-sms-register ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /app/modem-sms-register .
COPY config.example.json ./
RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app

USER appuser

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep modem-sms-register || exit 1

CMD ["./modem-sms-register"]