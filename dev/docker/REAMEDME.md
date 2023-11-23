
### build-env
```shell

docker build -f ./Dockerfile_x86 -t build-env-x86 .
docker tag eb6f8777e3c3 registry.cn-hongkong.aliyuncs.com/fbc/build-env-x86:lastest
docker push registry.cn-hongkong.aliyuncs.com/fbc/build-env-x86:lastest

```


### build-fullNode-mainnet
```shell

git https://github.com/fibonacci-chain/fbc/fbc.git
cd ./fbc
git checkout v1.6.8.6
cd ..
docker build -t fullnode-mainnet:v1.6.8.6 .


docker login --username=zjf@1633522246959130 registry.cn-hongkong.aliyuncs.com
docker tag [ImageId] registry.cn-hongkong.aliyuncs.com/fbc/fullnode:[镜像版本号]
docker push registry.cn-hongkong.aliyuncs.com/fbc/fullnode:[镜像版本号]


```


### delete image
```shell
# 停止docker
docker stop $(docker ps -a | grep "Exited" | awk '{print $1 }')
# 删除docker
docker rm $(docker ps -a | grep "Exited" | awk '{print $1 }')
# 删除images
docker rmi $(docker images | grep "none" | awk '{print $3}')


```

### stop docker
```shell
#两个命令均可使用，二选一即可
kill -2 ${pid}
kill -15 ${pid}

docker stop -t 1200  或者 docker-compose down -t 1200  # 停止容器

```


```shell

init
sudo apt install docker.io 
sudo apt install docker-compose 
sudo docker login --username=zjf@1633522246959130 registry.cn-hongkong.aliyuncs.com 


sudo vi docker-compose.yml 
sudo docker-compose pull 

fbchaind stop

sudo docker-compose up -d 
sudo docker-compose logs -f --tail 100 


```


```shell

OLD 机器:

x86 goleveldb repair-state

sudo ./cli/fbchaind repair-state --start-height 10925700 --db_backend goleveldb --home ~/.fbchaind


```