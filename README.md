# FBC


#### Build docker
```shell
docker build -t fullnode-mainnet:v1.6.8.5 .

docker tag {imageId} {registryUrl}/fbc/fullnode:v1.6.8.5

docker push {registryUrl}/fbc/fullnode:v1.6.8.5

/dev/docker

sudo docker-compose up -d 

```

### Join Fibonacci Mainnet

##### genesis file
```shell

build/genesis.json

```

##### public seed node

```shell

757fe84d81e1d09fb24c76265a02f36b76c2bc5b@16.162.64.131:26656
c176d3003ca8b6f66cebfd2df091a3a5c5c37113@16.163.133.239:26656
9bffbd7b5b22dc2e05cc191ad964949bbcee6751@43.198.116.254:26656

```

##### run a node
```shell
make mainnet WITH_ROCKSDB=true

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

```

##### run a node env options
```shell

============================Run Validator/RPC ENV==================================

      - FBC_LOG_FILE=/root/logs/fbchaind.log
      - FBC_LOG_STDOUT=false
      - FBC_DB_BACKEND=goleveldb
      - FBC_ELAPSED=DeliverTxs=2,Round=1,CommitRound=1,Produce=1
      - FBC_LOG_LEVEL=main:info,iavl:info,*:error,tx-receiver:info
      - FBC_NODE_MODE=val #archive, val, rpc
      - FBC_FAST_QUERY=true
      - FBC_DEBUG_API=false
      - FBC_MIN_GAS_PRICE=0.0000001fibo
      - FBC_ADDR_BOOK_STRICT=false
      - FBC_UNSAFE_CORS=true
      
============================Run ARCHIVE ENV========================================

      - FBC_LOG_FILE=/root/logs/fbchaind.log
      - FBC_LOG_STDOUT=false
      - FBC_DB_BACKEND=rocksdb
      - FBC_ELAPSED=DeliverTxs=2,Round=1,CommitRound=1,Produce=1
      - FBC_LOG_LEVEL=main:info,iavl:info,*:error,tx-receiver:info
      - FBC_NODE_MODE=archive #archive, val, rpc
      - FBC_FAST_QUERY=false
      - FBC_DEBUG_API=false
      - FBC_MIN_GAS_PRICE=0.0000001fibo
      - FBC_ADDR_BOOK_STRICT=false
      - FBC_IAVL_ENABLE_ASYNC_COMMIT=false
      - FBC_ROCKSDB_OPTS=max_open_files=20000
      - FBC_IAVL_CACHE_SIZE=5000000
      - FBC_IAVL_FAST_STORAGE_CACHE_SIZE=5000000

```
