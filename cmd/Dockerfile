# Build from repository root:
#   docker compose build
#
# Run (reads .env automatically):
#   docker compose up

FROM golang:1.26.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/server ./server

USER app

EXPOSE 80

CMD ["./server"]
