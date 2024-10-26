FROM golang:1.23 AS build

WORKDIR /tmp
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest && \
    CGO_ENABLED=0 go install -tags 'postgres sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go install ./...

FROM alpine:latest

COPY migrations /srv/migrations
COPY --from=build /go/bin /srv

ENTRYPOINT ["/srv/pulsar-alice"]