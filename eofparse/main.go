package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
			var codes []string
			for _, code := range c.Code {
				codes = append(codes, fmt.Sprintf("%x", code))
			}
			fmt.Printf("OK %v\n", strings.Join(codes, ","))
		} else {
			fmt.Printf("OK 0x\n")
		}
	}
}
