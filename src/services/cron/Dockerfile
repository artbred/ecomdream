FROM golang:1.19 AS builder

LABEL version="1.0"

RUN mkdir /go/src/ecomdream
COPY . /go/src/ecomdream
WORKDIR /go/src/ecomdream

RUN go mod download
RUN cd src/services/cron && CGO_ENABLED=0 GOOS=linux go build

FROM scratch AS production

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /go/src/ecomdream .

CMD ["./src/services/cron/cron"]
