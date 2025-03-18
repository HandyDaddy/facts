# Используем последнюю версию образа Alpine Linux с Go
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/app /app/
COPY config/ /app/config/

WORKDIR /app
ENV BEARERTOKEN ${BEARERTOKEN}

CMD ["./app"]