# Build container
FROM golang:alpine AS builder

RUN apk update && apk upgrade && apk add git zlib-dev gcc musl-dev

COPY . /go/src/github.com/Miniplays-Tickets/dashboard
WORKDIR /go/src/github.com/Miniplays-Tickets/dashboard

RUN git submodule update --init --recursive --remote

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -tags=jsoniter \
    -trimpath \
    -o main cmd/api/main.go

# Prod container
FROM alpine:latest

RUN apk update && apk upgrade && apk add curl

COPY --from=builder /go/src/github.com/Miniplays-Tickets/dashboard/locale /srv/dashboard/locale
COPY --from=builder /go/src/github.com/Miniplays-Tickets/dashboard/main /srv/dashboard/main

RUN chmod +x /srv/dashboard/main

RUN adduser container --disabled-password --no-create-home
USER container
WORKDIR /srv/dashboard

CMD ["/srv/dashboard/main"]