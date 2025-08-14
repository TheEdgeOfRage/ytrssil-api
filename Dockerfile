FROM golang:1.24-trixie AS builder

COPY go.mod go.sum /app/
WORKDIR /app/
RUN go mod download

COPY . /app/
RUN go build -o dist/ytrssil-api cmd/main.go

FROM debian:trixie-slim AS api
RUN apt update \
	&& apt install -y ca-certificates curl \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/dist/ /app/
ENTRYPOINT ["/app/ytrssil-api"]

FROM migrate/migrate AS migrations
COPY ./migrations/ /migrations/
