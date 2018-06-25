# Build Glbchain-dev in a stock Go builder container
FROM golang:1.9-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go-lbchain-devereum
RUN cd /go-lbchain-devereum && make glbchain-dev

# Pull Glbchain-dev into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-lbchain-devereum/build/bin/glbchain-dev /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp 30304/udp
ENTRYPOINT ["glbchain-dev"]
