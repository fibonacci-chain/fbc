FROM build-env as builder
WORKDIR /root
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct

COPY fbc ./fbc

RUN cd fbc && \
    make mainnet WITH_ROCKSDB=true LINK_STATICALLY=true

FROM golang:1.17.2-alpine
RUN apk add --no-cache bash
WORKDIR /root
COPY --from=builder $GOPATH/bin/fbchaind $GOPATH/bin/fbchaind
COPY --from=builder $GOPATH/bin/fbchaincli $GOPATH/bin/fbchaincli
COPY start.sh .
RUN chmod +x start.sh
CMD ["sh", "-c", "/root/start.sh"]
EXPOSE 26656 26657 26660 6060 8545 8546