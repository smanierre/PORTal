FROM node:20.18 AS frontend
WORKDIR /ui
COPY ui/package.json /ui/package.json
COPY ui/package-lock.json /ui/package-lock.json

RUN npm install

COPY ui/ /ui

RUN npm run build

FROM golang:1.22.8-bookworm AS backend

WORKDIR /build

COPY go.mod /build/go.mod
COPY go.sum /build/go.sum
RUN go mod download

COPY . /build

RUN go build -o PORTal main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=backend /build/PORTal /app/PORTal
COPY --from=frontend /ui/dist /app/dist

EXPOSE 8080
CMD [ "/app/PORTal" ]