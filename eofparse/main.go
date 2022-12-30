package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func main() {
	work()
}

func work() {
	// The input is assumed to be an EOF1 container verified against Shanghain instructionset.
	jt := vm.NewShanghaiEOFInstructionSetForTesting()
	var c vm.Container
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		blob := common.FromHex(scanner.Text())
		err := c.UnmarshalBinary(blob)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		err = c.ValidateCode(&jt)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		if len(c.Code) > 0 {
			fmt.Printf("OK %x\n", c.Code[0])
		} else {
			fmt.Printf("OK 0x\n")
		}
	}
}
