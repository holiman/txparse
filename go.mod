module github.com/holiman/txparse

go 1.19

require (
	github.com/ethereum/go-ethereum v1.10.26
	github.com/holiman/goevmlab v0.0.0-20221214131815-ba1151503917
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/holiman/big v0.0.0-20221017200358-a027dc42d04e // indirect
	github.com/holiman/uint256 v1.2.1 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
)

//replace github.com/ethereum/go-ethereum => github.com/lightclient/go-ethereum b30a56bf4a9713bfe183b690866a57934e69ec2e
replace github.com/ethereum/go-ethereum => /home/user/go/src/github.com/ethereum/go-ethereum
