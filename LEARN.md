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

#### lease 租期

COMMANDS:

| 子命令     | 描述                       | 描述                                |
|------------|----------------------------|-------------------------------------|
| grant      | Creates leases             | 创建租期并返回一个租期的 ID         |
| keep-alive | Keeps leases alive (renew) | 对指定租期 ID 进行包活              |
| list       | List all active leases     | 获取所有有效的 lease                |
| revoke     | Revokes leases             | 删除租期,与之相关联的键都会被删除掉 |
| timetolive | Get lease information      | 根据 ID 获取指定 lease 的信息       |

```bash
$etcdctl lease grant 1000 # 申请一个租期
lease 32696eed79ac2a0a granted with TTL(1000s) # 申请成功,得到一个租期 ID
$etcdctl lease timetolive 32696eed79ac2a0a # 查看租期的信息
lease 32696eed79ac2a0a granted with TTL(1000s), remaining(974s)
$etcdctl put test_lease some_value --lease 32696eed79ac2a0a # put 一个键 并挂载上面申请的租期
OK
$etcdctl put test_lease1 some_value --lease 32696eed79ac2a0a # put 第二个键 并挂载上面申请的同一个租期
OK
$etcdctl put test_lease2 some_value --lease 32696eed79ac2a0a # put 第三个键 并挂载上面申请的同一个租期
OK
$etcdctl lease timetolive --keys  32696eed79ac2a0a # 获取租期的信息并使用`--keys`指定获取挂载该租期的所有键
lease 32696eed79ac2a0a granted with TTL(1000s), remaining(885s), attached keys([test_lease test_lease1 test_lease2])
$etcdctl get --from-keys ""
test_lease
test_lease1
test_lease2
$etcdctl lease keep-alive 32696eed79ac2a0a # 对租期进行保活操作,每次刷新租期的到期时间以保证租期一直存在
lease 32696eed79ac2a0a keepalived with TTL(1000)
...
lease 32696eed79ac2a0a keepalived with TTL(1000)
$etcdctl lease revoke 32696eed79ac2a0a # 删除租期,挂在在该租期的所有键值都会被删除
lease 32696eed79ac2a0a revoked
$etcdctl get ''  --from-keys # 删除完租期之后,挂载在该租期上的所有键都会被删除
```

