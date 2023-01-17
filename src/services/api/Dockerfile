FROM golang:1.19 AS builder

LABEL version="1.0"

RUN mkdir /go/src/api
COPY . /go/src/api
WORKDIR /go/src/api

RUN go mod download
RUN cd src/services/api && CGO_ENABLED=0 GOOS=linux go build

FROM scratch AS production

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /go/src/api .

CMD ["./src/services/api/api"]
