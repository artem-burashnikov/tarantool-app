package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  "tarantool-storage:3301",
		User:     "storage",
		Password: "sesame",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return
	}

	fmt.Println("Connection established")

	time.Sleep(5 * time.Second)

	conn.CloseGraceful()
	fmt.Println("Connection is closed")
}
