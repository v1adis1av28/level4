FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /app/demoservice ./cmd

FROM alpine:latest
WORKDIR /app

RUN mkdir -p /app/config

COPY --from=builder /app/demoservice .
COPY config/dev.yml /app/config/dev.yml

EXPOSE 8080
EXPOSE 6060
CMD ["./demoservice"]