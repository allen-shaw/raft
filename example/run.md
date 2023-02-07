# How To Run

## 1. 管理命令
```shell
## 启动node1
./raftexample --id 1 --cluster http://127.0.0.1:12379 --port 12380
## node2加入
curl -L http://127.0.0.1:12380/2 -XPOST -d http://127.0.0.1:22379
./raftexample --id 2 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379 --port 22380 --join
## node3加入
curl -L http://127.0.0.1:12380/3 -XPOST -d http://127.0.0.1:32379
./raftexample --id 3 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port 32380 --join
```

## 2. 数据操作命令
```shell
## set kv
curl -L http://127.0.0.1:22380/my-key -XPUT -d bar
## get kv
curl -L http://127.0.0.1:22380/my-key
```