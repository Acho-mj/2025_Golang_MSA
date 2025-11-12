# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY backend backend
COPY proto proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /workspace/bin/order-service ./backend/services/order

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /workspace/bin/order-service /app/order-service

USER 65532:65532

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/order-service"]

