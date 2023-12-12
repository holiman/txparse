package main

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"

	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

func main() {
	for {
		fmt.Fprintf(os.Stdout, "%#x\n", fuzz())
		//break
	}
}

func rndBool(a int) bool {
	return rand.Intn(a) == 0
}

func rndBytes(min, max int) []byte {
	var size = min
	if max > min {
		for i := 0; i < 10; i++ {
			size += min + rand.Intn(max-min)
		}
		size = size / 10
	}
	buf := make([]byte, size)
	crand.Read(buf)
	return buf
}

func fuzz() []byte {
	var lists []int
	b := rlp.NewEncoderBuffer(nil)
	lists = append(lists, b.List())
	rounds := 0
	for {

		if rndBool(2) {
			b.WriteBigInt(big.NewInt(1))
		} else if rndBool(2) {
			b.WriteUint64(rand.Uint64())
		} else if rndBool(2) {
			b.WriteBytes(rndBytes(32, 32))
		} else {
			b.WriteBytes(rndBytes(20, 20))
		}
		// start a list
		if rndBool(20) {
			l := b.List()
			//fmt.Fprintf(os.Stderr, "starting list %d\n", l)
			lists = append(lists, l)
		}
		rounds++
		if rounds >= 15 {
			break
		}

		if rndBool(15 - rounds) {
			break
		}

	}
	// end all lists
	for i := len(lists) - 1; i >= 0; i-- {
		//fmt.Fprintf(os.Stderr, "ending list %d\n", lists[i])

		b.ListEnd(lists[i])
	}
	if rndBool(5) { // 1 in 5
		return b.ToBytes()
	}
	if true || rndBool(2) { // 1 in 2
		var buf = []byte{byte(1 + rand.Int()%3)}
		buf = append(buf, b.ToBytes()...)
		return buf
	}
	if rndBool(2) || rndBool(2) { // 1 in 2 or 1 in 2
		var buf = make([]byte, 1)
		crand.Read(buf)
		buf = append(buf, b.ToBytes()...)
		return buf
	}
	var buf = make([]byte, 2)
	crand.Read(buf)
	buf = append(buf, b.ToBytes()...)
	return buf
}
