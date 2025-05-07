package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	onlyGood = flag.Bool("onlygood", false, "used to filter 'only good' inputs")
	onlyBad  = flag.Bool("onlybad", false, "used to filter 'only bad' inputs")
	dump     = flag.Bool("txdump", false, "show detailed transaction dump")
)

func main() {
	flag.Parse()
	var (
		// The input is assumed to be Prague-enabled mainnet (chainid=1) transaction.
		signer   = types.NewPragueSigner(new(big.Int).SetInt64(1))
		scanner  = bufio.NewScanner(os.Stdin)
		toRemove = regexp.MustCompile(`[^0-9A-Za-z]`)
	)
	scanner.Buffer(make([]byte, 1024*1024), 5*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		sanitized := toRemove.ReplaceAll(line, []byte{})
		data := common.FromHex(string(sanitized))
		sender, err := parseSender(signer, data, *dump)
		if err != nil {
			if *onlyGood {
				continue // ignore the input
			}
			if *onlyBad {
				fmt.Fprintln(os.Stderr, string(line))
			}
			fmt.Printf("err: %v\n", err)
			continue
		}
		if *onlyBad {
			continue
		}
		fmt.Printf("%#x\n", sender)
		if *onlyGood {
			fmt.Fprintln(os.Stderr, string(line))
		}
	}
}

func parseSender(signer types.Signer, data []byte, dump bool) (common.Address, error) {
	tx := new(types.Transaction)

	if err := tx.UnmarshalBinary(data); err != nil {
		return common.Address{}, err
	}
	if dump {
		d, _ := json.MarshalIndent(tx, "  ", "  ")
		fmt.Fprintln(os.Stderr, string(d))
	}
	if err := extendedValidation(tx); err != nil {
		return common.Address{}, err
	}
	return signer.Sender(tx)
}

// extendedValidation is validation that is normally not performed during RLP-decoding,
// but instead happens at a later stage.
func extendedValidation(tx *types.Transaction) error {

	// state_transition.go:318
	//if hashes := tx.BlobHashes(); len(hashes) == 0 {
	//	return fmt.Errorf("blobless blob transaction")
	//}
	//if len(hashes) > params.MaxBlobGasPerBlock/params.BlobTxBlobGasPerBlob {
	//	return fmt.Errorf("too many blobs in transaction: have %d, permitted %d", len(hashes), params.MaxBlobGasPerBlock/params.BlobTxBlobGasPerBlob)
	//}
	if tx.BlobTxSidecar() != nil {
		return errors.New("blob-tx with sidecar (not consensus-encoding)")
	}
	// The geth state transition does not explicitly check whether the value is below
	// 256 bits in size. This is implicitly validated by checking that the sender
	// has sufficient balance to cover the tx cost.
	if tx.Value().BitLen() > 256 {
		return errors.New("value larger than 256 bits")
	}
	return nil
}
