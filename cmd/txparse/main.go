package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

func main() {
	var (
		// The input is assumed to be Cancun-enabled mainnet (chainid=1) transaction.
		signer   = types.NewCancunSigner(new(big.Int).SetInt64(1))
		scanner  = bufio.NewScanner(os.Stdin)
		toRemove = regexp.MustCompile(`[^0-9A-Za-z]`)
	)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		sanitized := toRemove.ReplaceAllString(line, "")
		data := common.FromHex(sanitized)
		sender, err := parseSender(signer, data)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}
		fmt.Printf("%#x\n", sender)
	}
}

func parseSender(signer types.Signer, data []byte) (common.Address, error) {
	tx := new(types.Transaction)

	if err := tx.UnmarshalBinary(data); err != nil {
		return common.Address{}, err
	}
	return signer.Sender(tx)
}
