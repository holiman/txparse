package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

func Fuzz(f *testing.F) {

	//	binaries := "/home/user/go/src/github.com/holiman/txparse/eofparse/eofparse,/home/user/go/src/github.com/holiman/txparse/eofparse/eofparse"
	binaries := os.ReadFile("binaries.txt")
	bins := strings.Split(binaries, ",")
	if len(binaries) < 2 {
		fmt.Printf("Usage: comparer parser1,parser2,... \n")
		fmt.Printf("Pipe input to process")
		f.Fatal("error")
	}
	var inputs = make(chan string)
	go func() {
		doit(bins, inputs)
	}()
	fil, err := os.Open("../eofparse/all.input")
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

	f.Fuzz(func(t *testing.T, data []byte) {
		_ = testUnmarshal(data) // This is for coverage guidance
		inputs <- fmt.Sprintf("%#x", data)
	})

}

var jt = vm.NewShanghaiEOFInstructionSetForTesting()

func testUnmarshal(blob []byte) string {
	var c vm.Container
	if err := c.UnmarshalBinary(blob); err != nil {
		return fmt.Sprintf("err: %v\n", err)
	}
	if err := c.ValidateCode(&jt); err != nil {
		return fmt.Sprintf("err: %v\n", err)
	}
	if len(c.Code) == 0 {
		return "OK Ox\n"
	}
	var codes []string
	for _, code := range c.Code {
		codes = append(codes, fmt.Sprintf("%x", code))
	}
	return fmt.Sprintf("OK %v\n", strings.Join(codes, ","))
}
