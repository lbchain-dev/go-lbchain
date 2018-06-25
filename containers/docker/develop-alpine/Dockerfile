FROM alpine:3.7

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 https://github.com/lbchain-devchain/go-lbchain-dev && \
  (cd go-lbchain-devereum && make glbchain-dev) && \
  cp go-lbchain-devereum/build/bin/glbchain-dev /glbchain-dev && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /go-lbchain-devereum && rm -rf /var/cache/apk/*

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/glbchain-dev"]
