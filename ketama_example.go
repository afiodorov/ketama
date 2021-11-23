// Copyright 2016 Chao Wang <hit9@icloud.com>

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"github.com/afiodorov/ketama"
)

func main() {
	ring := ketama.NewRing([]*ketama.Node{
		ketama.NewNode("127.0.0.1:8000", "binding data0", 1),
		ketama.NewNode("127.0.0.1:8001", "binding data1", 1),
		ketama.NewNode("127.0.0.1:8002", "binding data2", 1),
		ketama.NewNode("127.0.0.1:8003", "binding data3", 1),
		ketama.NewNode("127.0.0.1:8004", "binding data3", 1),
	})
	fmt.Printf("Get a server by key \"key1\": %v\n", ring.Get("key1"))
	fmt.Printf("Get a server by key \"key2\": %v\n", ring.Get("key2"))
	fmt.Printf("Get a server by key \"key3\": %v\n", ring.Get("key3"))
	fmt.Printf("Get a server by key \"key4\": %v\n", ring.Get("key4"))
	fmt.Printf("Get a server by key \"key5\": %v\n", ring.Get("key5"))
	fmt.Printf("Get a server by key \"key1\" again: %v\n", ring.Get("key1"))

	fmt.Printf("Get a server by key \"key1\" after assigned node goes down: %v\n", ring.GetIgnoringFailed("key1", map[string]struct{}{
		"127.0.0.1:8000": struct{}{},
	}))

	fmt.Printf("Get a server by key \"key1\" after all nodes are down: %v\n", ring.GetIgnoringFailed("key1", map[string]struct{}{
		"127.0.0.1:8000": struct{}{},
		"127.0.0.1:8001": struct{}{},
		"127.0.0.1:8002": struct{}{},
		"127.0.0.1:8003": struct{}{},
		"127.0.0.1:8004": struct{}{},
	}))
}
