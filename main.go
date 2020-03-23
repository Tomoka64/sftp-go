package main

import (
	"context"
	"fmt"

	"github.com/Tomoka64/sftp-go-sample/csv"
	"github.com/Tomoka64/sftp-go-sample/sftp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := newConnectionThroughProxy(ctx)
	if err != nil {
		panic(err)
	}

	cli, err := sftp.NewClient(conn)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	testData := "./test_data/test.csv"
	c := csv.New()
	records, err := c.Read(testData)
	if err != nil {
		panic(err)
	}
	bs, err := c.Bytes(records)
	if err != nil {
		panic(err)
	}

	dest := "test.csv"
	err = cli.Upload(dest, bs)
	if err != nil {
		panic(err)
	}

	if !cli.Exist(dest) {
		fmt.Println("FAIL")
	}
	fmt.Println("SUCCESSFULLY uploaded file")
}
