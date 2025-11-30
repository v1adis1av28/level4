FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /app/producer ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/producer .

#EXPOSE 8080
CMD ["./producer"]