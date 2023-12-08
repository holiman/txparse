package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	var (
		// The input is assumed to be Cancun-enabled mainnet (chainid=1) transaction.
		signer   = types.NewCancunSigner(new(big.Int).SetInt64(1))
		scanner  = bufio.NewScanner(os.Stdin)
		toRemove = regexp.MustCompile(`[^0-9A-Za-z]`)
	)
	scanner.Buffer(make([]byte, 1024*1024), 5*1024*1024)
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
	return nil
}
