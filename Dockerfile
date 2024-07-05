FROM golang:1.22 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o podcast-sync-server

FROM debian:bookworm
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends --assume-yes libsqlite3-0

COPY --from=builder /app/podcast-sync-server /app/podcast-sync-server

CMD ["/app/podcast-sync-server", "serve"]
