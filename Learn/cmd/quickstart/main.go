package main

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
)

var client *clientv3.Client

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		log.Fatalln("failed in new client", err)
	}
}

func main() {
	setKey()
	getWithRange()
}

func GrantLeaseAndKeepAlive() clientv3.LeaseID {
	lease, err := client.Grant(context.TODO(), 10)
	if err != nil {
		log.Fatalln("failed in grant lease", err)
	}
	log.Println(lease.ID)
	reponseChan, err := client.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		log.Fatalln("failed in keep alive", err)
	}
	go func() {
		for response := range reponseChan {
			log.Println("get keepalive reponse:", response)
		}
	}()
	return lease.ID
}

func PutWithLease(id clientv3.LeaseID) {
	put, err := client.Put(context.TODO(), "bar", "foo", clientv3.WithLease(id))
	if err != nil {
		log.Fatalln("failed in put", err)
	}
	log.Println("put response:", put)
}

func setKey() {
	keys := []string{
		"bar1",
		"bar2",
		"bar3",
	}
	for _, key := range keys {
		_, err := client.Put(context.TODO(), key, "foo")
		if err != nil {
			log.Fatalln("failed in put:", err)
		}
	}
}

func getWithPrefix() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	response, err := client.Get(ctx, "bar", clientv3.WithPrefix())
	if err != nil {
		log.Fatalln("failed in get:", err)
	}
	log.Println("get with prefix response:")
	for _, kv := range response.Kvs {
		log.Println(kv)
	}
}

func getWithRange() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//注意这个区间是 [bar1,bar4) 不包括 bar4
	response, err := client.Get(ctx, "bar1", clientv3.WithRange("bar4"))
	if err != nil {
		log.Fatalln("failed in get range:", err)
	}
	for _, kv := range response.Kvs {
		log.Println(kv)
	}
}

func revokeLease(id clientv3.LeaseID) {
	_, err := client.Revoke(context.TODO(), id)
	if err != nil {
		log.Fatalln("failed in revoke:", err)
	}
}

func errorNotNil(err error, msg string) {
	if err != nil {
		log.Fatalf("failed in %s:%v", msg, err)
	}
}

func Set() {
	leases, err := client.Leases(context.Background())
	errorNotNil(err, "failed in  list leases")
	for lease := range leases.Leases {
		log.Println(lease)
	}
}
