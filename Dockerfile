FROM golang:1.22.8-bookworm AS builder

WORKDIR /build
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
RUN go mod download

COPY . /build

RUN go build -o PORTal main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /build/PORTal /app/PORTal

EXPOSE 8080
CMD [ "/app/PORTal" ]
