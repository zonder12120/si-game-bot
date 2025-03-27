# STAGE 1 BUILD

FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o tg-bot .

# STAGE 2 RUN

FROM alpine:3.18

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/tg-bot .

RUN chown -R appuser:appgroup /app

USER appuser

ENTRYPOINT ["/app/tg-bot"]