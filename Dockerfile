# 1
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.su[m] .
RUN go mod download

COPY . .
RUN go build -o main ./cmd/ssh/

# 2
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3737

CMD ["./main"]
