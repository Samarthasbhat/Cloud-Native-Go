package main

//Maps are not atomic in Go. So we need to use mutex.

import (

"sync"

)

// By using the mutex in this way, we can ensure that exactly one process has exclusive access to a particular resource.
// We can use sync package for that (RWmutex)

var myMap = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

myMap.Lock()  // Take a write Lock
myMap.m["some key"] = "some_value"
myMap.Unlock() // Release write Lock

// Read lock

myMap.RLock()
value := myMap.m["some_key"]
myMap.RUnlock()

fmt.Println("some_key:", value)

// Read locks are less restrictive than write locks 
// Now we can implement the RWMutex to Simple APIdocker run --name pg -e POSTGRES_USER=test -e POSTGRES_PASSWORD=hunter2 \
  -v "/c/Users/samar/OneDrive/Documents/Cloud Native Go/Building Cloud Native Service/storing-state-db/pg_hba.conf":/etc/postgresql/conf.d/pg_hba.conf:ro \
  -p 5432:5432 -d postgres

