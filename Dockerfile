FROM golang:1.16-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o bin .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin /app/bin

COPY . .

RUN ls /app

RUN chmod +x /app/bin

ENTRYPOINT ["/app/bin"]
