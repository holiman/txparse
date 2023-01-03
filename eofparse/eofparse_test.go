package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

/*
*
Corpus winds up e.g. here:

	[user@work txparse]$ cat ~/.cache/go-build/fuzz/github.com/holiman/txparse/Fuzz/6293d045c6cf3a446e7463c0"EF71eab95457c914513223000e097a0823c82dd4
	go test fuzz v1
	[]byte("\x01\xf8A000\x8200\xb4000000000000000000000000000000000000000000000000000000000000000000000000000000000")

It can also be overridden:

	$ GOCACHE=`pwd`/gen_corpus  go test . -fuzz Fuzz
*/
var jt = vm.NewShanghaiEOFInstructionSetForTesting()

func testUnmarshal(t *testing.T, data []byte) {
	var c vm.Container

	err := c.UnmarshalBinary(data)
	if err != nil {
		return
	}
	err = c.ValidateCode(&jt)
	if err != nil {
		return
	}
	out := c.MarshalBinary()
	if !bytes.Equal(out, data) {
		panic(fmt.Sprintf("Marshal/Unmarshal mismatch. \nInput:  %#x\nOutput: %#x\n", data, out))
	}
}

func Fuzz(f *testing.F) {
	fil, err := os.Open("all.input")
	if err != nil {
		f.Fatal(err)
	}
	defer fil.Close()
	corpi := 0
	scanner := bufio.NewScanner(fil)
	for scanner.Scan() {
		input := common.FromHex(scanner.Text())
		if len(input) > 0 {
			f.Add(input)
			corpi++
		}
	}
	f.Logf("Added seed corpus, %d items\n", corpi)
	f.Fuzz(testUnmarshal)
}
