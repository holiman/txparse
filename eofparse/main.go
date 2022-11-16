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

	jt := vm.NewLatestInstructionSetForTesting()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		blob := common.FromHex(scanner.Text())
		c, err := vm.ParseEOF1Container(blob)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		if err := c.ValidateCode(&jt); err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		fmt.Printf("OK %x\n", c.CodeAt(0))
	}
}
