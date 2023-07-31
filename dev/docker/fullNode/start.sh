#!/bin/bash

if [ ! -d "/root/.fbchaind/config" ]; then
    fbchaind init fullnode --chain-id fbc-1230
    wget https://fibochain.s3-ap-east-1.amazonaws.com/0/ayOqRj7bk6dMmGp9j38jtukPS_genesis.json -O /root/.fbchaind/config/genesis.json
fi

mkdir /cli
cp $GOPATH/bin/fbchaind /cli/fbchaind
cp $GOPATH/bin/fbchaincli /cli/fbchaincli

fbchaind start \
    --chain-id fbc-1230 \
    --unsafe-cors=${FBC_UNSAFE_CORS:-"false"} \
    --rest.laddr tcp://0.0.0.0:8545 \
    --log_level ${FBC_LOG_LEVEL} \
    --node-mode ${FBC_NODE_MODE} \
    --fast-query=${FBC_FAST_QUERY} \
    --debug-api=${FBC_DEBUG_API} \
    --minimum-gas-prices ${FBC_MIN_GAS_PRICE} \
    --p2p.persistent_peers=${FBC_PERSISTENT_PEERS} \
    --p2p.addr_book_strict=${FBC_ADDR_BOOK_STRICT} \
    --iavl-enable-async-commit=${FBC_IAVL_ENABLE_ASYNC_COMMIT:-"true"} \
    --iavl-cache-size=${FBC_IAVL_CACHE_SIZE:-10000000} \
    --iavl-fast-storage-cache-size=${FBC_IAVL_FAST_STORAGE_CACHE_SIZE:-10000000} \
    --rocksdb.opts=${FBC_ROCKSDB_OPTS:-""} \
    --db_backend ${FBC_DB_BACKEND}

#addr_book_strict: 使用严格的路由规则设置为true，使用本地网络或是私有网络设置为false，默认为true
#iavl-enable-async-commit 归档节点设置为false
#--iavl-enable-async-commit=false  #关闭iavl异步持久化
#--iavl-cache-size=xxx #iavl缓存，默认10000000，适当减小
#--iavl-fast-storage-cache-size=yyy  #fast storage缓存，默认10000000，适当减小
#--rocksdb.opts max_open_files=20000

#经调查，Rocksdb在归档节点占用内存大的主要原因是，为了保证Rocksdb的高读写性能，默认将rocksdb的参数设为max_open_files=-1。
#这导致启动Rocksdb时将DB中的所有file都打开并缓存，致使内存占用巨大
#启动节点时设置参数--rocksdb.opts max_open_files=20000
#参数值根据系统配置具体调整，该值设为-1表示不限制Rocksdb打开文件数；
#参数值>0时，值越大内存占用越大，DB读写性能越高；值越小内存占用越小，DB读写性能越低