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

1. `etcdctl`命令行界面,`etcdctl`支持读取环境变量作为参数.

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

因为 etcd 的日志都有版本(logIndex)，可以回到之前的版本(Index)，但是这样的代价就是牺牲了存储空间.当之前这些版本不再需要的时候,需要这些版本的最后的状态拍摄成快照并将这些版本删除.

`etcdctl compaction reversion`

#### watch 监测键值的修改

WATCH [options] [key or prefix] [range_end] [--] [exec-command arg1 arg2 ...]

当`watch`的键值有修改时会执行后面的命令,这个命令会持续运行直到出现错误或者被取消.像`get`命令一样,被`watch`的可以是一个键,一个键的范围,或者是所有包含指定前缀的键.
命令为:

- `watch`并打印修改内容:`etcdctl watch the_watched_key`
- `watch`并执行相应的命令:`etcdctl watch the_watched_key -- echo "the_watched_key is changed"`

#### 租期
