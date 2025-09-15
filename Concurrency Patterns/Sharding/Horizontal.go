package main

import (
	"fmt"
	"hash/fnv"
)


var servers = []string{"server1", "server2", "server3"}

func getShard(key string) string {
	h := fnv.New32a()
	h.Write([]byte(key))
	return servers[int(h.Sum32())%len(servers)]
}

func main() {
	users := []string{"alice", "bob", "charlie", "dave", "eve"}

	for _, user := range users {
		server := getShard(user)
		fmt.Printf("User %s is assigned to %s\n", user, server)
	}
}