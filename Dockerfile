FROM golang:1.22 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o podcast-sync-server

FROM debian:bookworm

LABEL org.opencontainers.image.source=https://github.com/whites11/podcast-sync-server
LABEL org.opencontainers.image.description="podcast-sync-server"
LABEL org.opencontainers.image.licenses=MIT

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes libsqlite3-0

COPY --from=builder /app/podcast-sync-server /app/podcast-sync-server

CMD ["/app/podcast-sync-server", "serve"]

