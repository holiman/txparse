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
		container, err := vm.ParseAndValidateEOF1Container(blob, &jt)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		fmt.Printf("OK %d\n", container.HeaderSize())
	}
}
