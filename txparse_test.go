package main

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func Fuzz(f *testing.F) {

	var (
		signer = types.NewLondonSigner(new(big.Int).SetInt64(1))
	)
	for _, tc := range []string{
		"0xf85f030182520894b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a801ca098ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4aa01887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a3",
		"0xf85f011082520894f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f001801ca06f0010ff4c31c2a6d0526c0d414e6cd01ad5d22e15bfff98af23867366b94d87a05413392d556119132da7056f8fb56a9138a36446a8a4ad7159c9d892d9f32284",
		"0xf8638080830f424094095e7baea6a6c7c4c2dfeb977efac326af552d87830186a0801ba0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0badf00d70ec28c94a3b55ec771bcbc70778d6ee0b51ca7ea9514594c861b1884",
		"0xf867078504a817c807830290409435353535353535353535353535353535353535358201578025a052f1a9b320cab38e5da8a8f97989383aab0a49165fc91c737310e4f7e9821021a052f1a9b320cab38e5da8a8f97989383aab0a49165fc91c737310e4f7e9821021",
		"0xf866068504a817c80683023e3894353535353535353535353535353535353535353581d88025a06455bf8ea6e7463a1046a0b52804526e119b4bf5136279614e0b1e8e296a4e2fa06455bf8ea6e7463a1046a0b52804526e119b4bf5136279614e0b1e8e296a4e2d	0x01f8bc018001826a4094095e7baea6a6c7c4c2dfeb977efac326af552d878080f85af858939e7baea6a6c7c4c2dfeb977efac326af552d87f842a00000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000180a05cbd172231fc0735e0fb994dd5b1a4939170a260b36f0427a8a80866b063b948a07c230f7f578dd61785c93361b9871c0706ebfa6d06e3f4491dc9558c5202ed36",
		"0x01f89a018001826a4094095e7baea6a6c7c4c2dfeb977efac326af552d878080f838f794a95e7baea6a6c7c4c2dfeb977efac326af552d87e1a0ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80a05cbd172231fc0735e0fb994dd5b1a4939170a260b36f0427a8a80866b063b948a07c230f7f578dd61785c93361b9871c0706ebfa6d06e3f4491dc9558c5202ed36",
		"0x01f87b018001826a4094095e7baea6a6c7c4c2dfeb977efac326af552d878080dad994a95e7baea6a6c7c4c2dfeb977efac326af552d87c382000180a05cbd172231fc0735e0fb994dd5b1a4939170a260b36f0427a8a80866b063b948a07c230f7f578dd61785c93361b9871c0706ebfa6d06e3f4491dc9558c5202ed36",
	} {
		f.Add(common.FromHex(tc)) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			f, err := os.OpenFile("./corpus", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				panic(err)
			}
			_, err = f.Write([]byte(fmt.Sprintf("0x%x\n", data)))
			if err1 := f.Close(); err1 != nil && err == nil {
				err = err1
			}
			if err != nil {
				panic(err)
			}
		}()

		tx := new(types.Transaction)
		if err := tx.UnmarshalBinary(data); err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		if sender, err := signer.Sender(tx); err != nil {
			fmt.Printf("err: %v\n", err)
		} else {
			fmt.Printf("%#x\n", sender)
		}
	})

}
