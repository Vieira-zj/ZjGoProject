package main

import (
	"fmt"
	"os"

	"src/demo.app/mongodb"
)

func main() {
	runN := 1

	// #1 mongo test
	if runN == 1 {
		mongodb.ConnectToDbAndTest()
	}

	// #2 bucket
	if runN == 2 {
		const (
			bucket = "test_bucket_transfer_data11"
			uid    = 1380469264
		)
		bucketInfo := mongodb.NewBucketInfo(bucket, uid)
		bucketInfo.QueryBucketInfo()
	}

	// #3 rs
	if runN == 3 {
		rsOp := mongodb.NewRsOperation()
		defer rsOp.Close()
		rsOp.InsertRsRecords()
		// rsOp.InsertRsRecordsParallel()
	}

	// #4 mongo op logs
	// cmd: ./main 10.200.30.11:8001
	if runN == 4 {
		mgoOp := mongodb.NewMgoOpertion(os.Args[1])
		defer mgoOp.Close()
		mgoOp.PrintMgoOpLogs()
	}

	fmt.Println("mongodb demo done.")
}
