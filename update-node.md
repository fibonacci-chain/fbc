### 1. Safely end the node process, otherwise it may damage the database.

```shell
# ps -ef | grep fbchaind

# sudo kill -2  ${pid}
```
### 2. Check if the node stops normally

```shell
# sudo docker-compose logs -f --tail 100
```

If exiting normally, the following log will be printed:

    fbchain    | E[2023-11-24|05:48:29.190][10] Stopping peer for error. module=p2p peer="Peer{MConn{43.198.41.15:26656} 5e05489a5ff1d61c8931ee2b7b35ef2a2ddb6e0f out}" err="read tcp 172.17.0.2:55224->43.198.41.15:26656: use of closed network connection"
    fbchain    | E[2023-11-24|05:48:29.192][10] Stopped accept routine, as transport is closed. module=p2p numPeers=0
    fbchain    | Close App
    fbchain    | I[2023-11-24|05:48:29.214][10] stopping iavl. module=iavl commitheight=16776857
    fbchain    | I[2023-11-24|05:48:29.215][10] CommitSchedule. module=iavl Height=16776857 Tree=acc IavlHeight=0 NodeNum=154 tpp=154 fss-add=12 fss-rm=0 trc="commitSchedule<0ms>, cacheNode<0ms>, Pruning<0ms>, batchSet-node<0ms>, batchSet-fss<0ms>, batchCommit<0ms>"
    fbchain    | I[2023-11-24|05:48:29.217][10] PruningSchedule. module=iavl Height=16773000 Tree=acc trc="pruningSchedule<2ms>, deleteVersion<1ms>, Commit<0ms>"
    fbchain    | I[2023-11-24|05:48:29.217][10] stopping iavl completed. module=iavl commitheight=16776857
    fbchain    | I[2023-11-24|05:48:29.220][10] stopping iavl. module=iavl commitheight=16776857
    fbchain    | I[2023-11-24|05:48:29.230][10] CommitSchedule. module=iavl Height=16776857 Tree=evm IavlHeight=0 NodeNum=1531 tpp=1531 fss-add=176 fss-rm=0 trc="commitSchedule<8ms>, cacheNode<1ms>, Pruning<0ms>, batchSet-node<2ms>, batchSet-fss<0ms>, batchCommit<4ms>"
    fbchain    | I[2023-11-24|05:48:29.256][10] PruningSchedule. module=iavl Height=16773000 Tree=evm trc="pruningSchedule<26ms>, deleteVersion<19ms>, Commit<6ms>"
    fbchain    | I[2023-11-24|05:48:29.256][10] stopping iavl completed. module=iavl commitheight=16776857
    fbchain    | I[2023-11-24|05:48:29.259][10] exiting.... module=main 
    fbchain exited with code 143

### To avoid unexpected situations where the node cannot be restored, it is recommended to rename the data directory ./fbchaind/data backup



### 4.  Modify the docker-compose.yml file

```shell

# sudo vi docker-compose.yml 

```
Note: This upgrade has modified the [image warehouse](https://hub.docker.com/search?q=eeebyte%2Ffullnode) and version number.

If it is a **arm architecture server**
- image: **eeebyte/fullnode-arm:v1.6.8.6**

If it is a **x86 architecture server**
- image: **eeebyte/fullnode-x86:v1.6.8.6**

The following configuration is for reference and needs to be configured according to your own path.
```yml
    version: "3"
    services:
      fbchain:
        container_name: fbchain
        #Please confirm if the server architecture running the node is ARM or x86
        image: eeebyte/fullnode-arm:v1.6.8.6
        #restart: unless-stopped
        environment:
          - FBC_LOG_FILE=/root/logs/fbchaind.log
          - FBC_LOG_STDOUT=false
          #Please confirm again whether the database used by the node is Rocksdb or Goleveldb, and modify this field according to your own database type
          - FBC_DB_BACKEND=rocksdb
          - FBC_ELAPSED=DeliverTxs=2,Round=1,CommitRound=1,Produce=1
          - FBC_LOG_LEVEL=main:info,iavl:info,*:error,tx-receiver:info
          - FBC_NODE_MODE=val #archive, val, rpc
          - FBC_FAST_QUERY=true
          - FBC_DEBUG_API=false
          - FBC_MIN_GAS_PRICE=0.0000001fibo
          - FBC_ADDR_BOOK_STRICT=false
          - FBC_UNSAFE_CORS=true
          - FBC_PERSISTENT_PEERS=
        volumes:
          - ./fbchaind:/root/.fbchaind
          - ./logs:/root/logs
          - ./cli:/cli
        ports:
          - "26656:26656"
          - "26657:26657"
          - "26660:26660"
          - "8546:8546"
          - "8545:8545"
          - "6060:6060"
```
### 5. Update image
```shell
sudo docker-compose pull 
```
### 6. Check the version number, the correct version number is v1.6.8.6
```
# sudo ./cli/fbchaind  version
v1.6.8.6
```
### 7. Starting node with a large amount of data may take several minutes to see the log results. Normal synchronization indicates successful updates.

```shell
# docker-compose up -d
# docker-compose logs -f --tail 100
```