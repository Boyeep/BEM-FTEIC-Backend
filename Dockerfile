FROM golang:1.25.12-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.22

RUN apk add --no-cache tzdata \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=builder /server /app/server

RUN mkdir -p /app/uploads && chown -R app:app /app

USER app

EXPOSE 8080

CMD ["/app/server"]
