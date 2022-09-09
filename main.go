package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	// The input is assumed to be London-enabled mainnet (chainid=1) transaction.
	var (
		signer  = types.NewLondonSigner(new(big.Int).SetInt64(1))
		scanner = bufio.NewScanner(os.Stdin)
		tx      = new(types.Transaction)
	)
	for scanner.Scan() {
		if err := tx.UnmarshalBinary(common.FromHex(scanner.Text())); err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		if sender, err := signer.Sender(tx); err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("%#x\n", sender)
		}
	}
}
