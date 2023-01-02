package main

import (
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/goevmlab/ops"
	"github.com/holiman/goevmlab/program"
)

func main() {
	work()
}

func work() {
	for {
		var c vm.Container
		numCodes := rand.Intn(3)
		for i := 0; i < numCodes; i++ {
			//code := make([]byte, 1+rand.Intn(50))
			//rand.Read(code)
			code, maxStack := genSimplProgam()
			c.Code = append(c.Code, code)
			var metadata = &vm.FunctionMetadata{
				Input:          uint8(rand.Intn(2)),
				Output:         uint8(rand.Intn(2)),
				MaxStackHeight: uint16(maxStack),
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

func genSimplProgam() ([]byte, int) {

	var p = program.NewProgram()
	var stackdepth = 0
	valids := ops.ValidOpcodes
	valids = append(valids, ops.OpCode(0x5c),
		ops.OpCode(0x5d),
		ops.OpCode(0x5e), ops.OpCode(0xb0), ops.OpCode(0xb1), ops.OpCode(0xb2))
	maxStack := 0
	//RJUMP    OpCode = 0x5c
	//RJUMPI   OpCode = 0x5d
	//RJUMPV OpCode = 0x5e
	//CALLF OpCode = 0xb0
	//RETF OpCode = 0xb1
	//JUMPF OpCode= 0xb2
	for nCases := 0; nCases < 40; nCases++ {
		op := ops.OpCode(valids[rand.Intn(len(valids))])
		pops := op.Pops()
		if pops != nil {
			if stackdepth < len(op.Pops()) {
				for i := 0; i < len(op.Pops()); i++ {
					p.Push(0)
					stackdepth++
				}
			}
			if stackdepth > maxStack {
				maxStack = stackdepth
			}
			stackdepth -= len(op.Pops())
			stackdepth += len(op.Pushes())
		} else {
			p.Push(0)
		}
		p.Op(op)
		if stackdepth > maxStack {
			maxStack = stackdepth
		}
	}
	p.Op(ops.STOP)
	return p.Bytecode(), maxStack
}
