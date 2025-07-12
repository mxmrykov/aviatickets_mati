FROM golang:1.23.0 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o aviatickets ./cmd/bot

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/aviatickets .
COPY --from=builder /app/configs configs/

EXPOSE 8080
CMD ["./aviatickets"]