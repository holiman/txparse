package main

import (
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/core/vm"
)

func main() {
	work()
}

func work() {
	for {
		var c vm.Container
		numCodes := rand.Intn(3)
		for i := 0; i < numCodes; i++ {
			code := make([]byte, 1+rand.Intn(50))
			rand.Read(code)
			c.Code = append(c.Code, code)
			var metadata = &vm.FunctionMetadata{
				Input:          uint8(rand.Intn(10)),
				Output:         uint8(rand.Intn(5)),
				MaxStackHeight: uint16(rand.Intn(10)),
			}
			if i == 0 {
				metadata.Input = 0
				metadata.Output = 0
			}
			c.Types = append(c.Types, metadata)
		}
		fmt.Printf("%x\n", c.MarshalBinary())
	}
}
