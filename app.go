package main

import (
	"enigmacamp.com/golatihanlagi/delivery"
	_ "github.com/lib/pq"
)

func main() {
	delivery.NewServer().Run()
}

// Unit Test
// 1. UseCase
// 2. Repository
// 3. Controller
