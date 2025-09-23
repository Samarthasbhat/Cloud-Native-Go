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
// Now we can implement the RWMutex to Simple API

