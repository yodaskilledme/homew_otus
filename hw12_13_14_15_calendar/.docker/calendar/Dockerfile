FROM golang:1.16-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o calendar cmd/calendar/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 8000

ENTRYPOINT ["./calendar"]