package main

import (
	"grpc/cmds"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmds.Execute()
}
