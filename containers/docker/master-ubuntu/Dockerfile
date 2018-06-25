FROM ubuntu:xenial

ENV PATH=/usr/lib/go-1.9/bin:$PATH

RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends golang-1.9 git make gcc libc-dev ca-certificates && \
  git clone --depth 1 --branch release/1.8 https://github.com/lbchain-devchain/go-lbchain-dev && \
  (cd go-lbchain-devereum && make glbchain-dev) && \
  cp go-lbchain-devereum/build/bin/glbchain-dev /glbchain-dev && \
  apt-get remove -y golang-1.9 git make gcc libc-dev && apt autoremove -y && apt-get clean && \
  rm -rf /go-lbchain-devereum

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/glbchain-dev"]
