### 1.  安全结束节点进程, 让节点正常退出，否则会损坏数据库。

```shell
# ps -ef | grep fbchaind

# sudo kill -2  ${pid}
```
### 2.  检查节点是否正常停止

```shell
# sudo docker-compose logs -f --tail 100
```

正常退出的话会打印如下日志：

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

### 3.  为了避免出现意外情况无法恢复节点运行，建议将数据目录 ./fbchaind/data 备份



### 4.  修改 docker-compose.yml 文件

```shell

# sudo vi docker-compose.yml 

```
注意：本次升级修改了[镜像仓库](https://hub.docker.com/search?q=eeebyte%2Ffullnode)及版本号
如果是**arm架构服务器**
- image: **eeebyte/fullnode-arm:v1.6.8.6**

如果是**x86架构服务器**
- image: **eeebyte/fullnode-x86:v1.6.8.6**

以下配置供参考使用，需要根据自己的路径去配置
```yml
    version: "3"
    services:
      fbchain:
        container_name: fbchain
        #请确认运行节点的服务器架构是arm还是x86
        image: eeebyte/fullnode-arm:v1.6.8.6
        #restart: unless-stopped
        environment:
          - FBC_LOG_FILE=/root/logs/fbchaind.log
          - FBC_LOG_STDOUT=false
          #请再三确认节点使用的数据库是rocksdb还是goleveldb，根据自己的数据库类型修改此字段
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
### 5. 更新镜像
```shell
sudo docker-compose pull 
```
### 6. 检查版本号,正确版本号是 v1.6.8.6
```
# sudo ./cli/fbchaind  version
v1.6.8.6
```
### 7. 启动节点，数据量比较大可能需要等几分钟才看到日志结果，正常同步即为更新成功。

```shell
# docker-compose up -d
# docker-compose logs -f --tail 100
```