package main

import (
	"fmt"
	redis "github.com/go-labx/example"
)

func main() {
	client, err := redis.NewClient(":6379")
	if err != nil {
		panic(err)
	}
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
		}
	}(client)

	_, err = client.Del("num")
	if err != nil {
		return
	}

	resp, err := client.Decrby("num", int64(100))
	fmt.Printf("resp = %v\nerr = %v", resp, err)

}
