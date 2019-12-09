# Learn

## 文档引用

- [master.Document](https://github.com/etcd-io/etcd/tree/master/Documentation)
- [mater.Document.Developing with etcd](https://github.com/etcd-io/etcd/tree/master/Documentation#developing-with-etcd)

### 开启本地集群

1. `go get github.com/mattn/goreman`获取 `goreman`
2. 运行 `goreman start`读取配置文件`./Procfile`开启本地测试集群

参考：

- [setting up local clusters](https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/local_cluster.md)

## 操作 etcd

1. `etcdctl`命令行界面

### K-V 操作

#### put key

- etcdctl put key value

#### 获取值

- range: etcdctl get foo foo100。获取 `[foo,foo100)`区间的值
- prefix: etcdctl get --prefix foo.获取所有以`foo`为前缀的key

#### 删除键

大致跟获取键值的语法一样，但是只是执行了删除操作

#### transaction 事物

事物从系统输入读取一系列`etcd`的请求将其作为原子事物一起执行，一个事物包含一些条件请求，一些请求用于在条件成立的时候运行，一些请求用于在条件不成立的时候运行。支持交互模式(`-i`)和非交互模式

```shell
./etcdctl txn <<<'mod("key") > "0"
put key1 "overwrote-key1"

put key1 "created-key1"
put key2 "some extra key"
'
```

#### compaction 日志压缩

因为 con