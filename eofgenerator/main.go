package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/goevmlab/fuzzing"
	"math/rand"
)

func main() {
	work()
}

func work() {
	for {
		var c vm.Container
		numCodes := 1024
		switch rand.Intn(5) {
		case 0:
			numCodes = 1
		case 1:
			numCodes = 2
		case 2:
			numCodes = 16
		case 3:
			numCodes = 1023
		default:
			numCodes = 1024
		}
		for i := 0; i < numCodes; i++ {
			code, maxStack := fuzzing.GenerateCallFProgram(numCodes)
			c.Code = append(c.Code, code)
			var metadata = &vm.FunctionMetadata{
				Input:          uint8(0),
				Output:         uint8(0),
				MaxStackHeight: uint16(maxStack),
			}
			if i == 0 {
				metadata.Input = 0
				metadata.Output = 0
			}
			c.Types = append(c.Types, metadata)
		}
		data := c.MarshalBinary()
		if len(data) > 0xc0f0 {
			data = data[:0xc000]
		}
		fmt.Printf("%x\n", data)
	}
}
