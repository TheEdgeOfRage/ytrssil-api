FROM golang:1.19-bullseye AS builder
RUN apt update && apt install -y make

# first copy just enough to pull all dependencies, to cache this layer
COPY go.mod go.sum Makefile /app/
WORKDIR /app/
RUN make setup

# lint, build, etc..
COPY . /app/
RUN make build

FROM debian:bullseye-slim
RUN apt update \
	&& apt install -y ca-certificates \
	&& apt clean \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/dist/ /app/
ENTRYPOINT ["/app/ytrssil-api"]
