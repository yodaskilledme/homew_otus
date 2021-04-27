FROM golang:1.16-alpine AS builder

WORKDIR /tmp/calendar

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/calendar -ldflags "$(LDFLAGS)" cmd/calendar/main.go

FROM alpine:latest

COPY --from=builder /tmp/calendar/bin/calendar /app/calendar

EXPOSE 8080

ENTRYPOINT ["./bin/calendar"]
