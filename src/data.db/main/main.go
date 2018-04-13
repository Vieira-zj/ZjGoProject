package main

import (
	"fmt"

	"data.db/memcached"
	"data.db/mongodb"
	"data.db/redis"
)

// build cmd: $ GOOS=linux GOARCH=amd64 go build
// $ scp main qboxserver@10.200.20.21:~/zhengjin/main
func main() {
	isMongodbTest := false
	if isMongodbTest {
		mongodb.ConnectToDbAndTest()
		// mongodb.InsertToRsDb()
		// mongodb.InsertToRsDbParallel()
		// cmd: ./main 10.200.30.11:8001
		// mongodb.PrintMongoOpLog()
	}

	isRedisTest := false
	if isRedisTest {
		// redis.ConnectToRedisAndTest()
		redis.MainRedis()
	}

	isMemTest := false
	if isMemTest {
		memcached.ConnectMemcacheAndTest()
	}

	fmt.Println("db main done.")
}