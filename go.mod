module github.com/holiman/txparse

go 1.19

require (
	github.com/ethereum/go-ethereum v1.10.26
	github.com/holiman/goevmlab v0.0.0-20230106175529-44c6800936c1
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VictoriaMetrics/fastcache v1.10.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/deckarep/golang-set/v2 v2.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/holiman/big v0.0.0-20221017200358-a027dc42d04e // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.2.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.5.0 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

//replace github.com/ethereum/go-ethereum => github.com/lightclient/go-ethereum b30a56bf4a9713bfe183b690866a57934e69ec2e
replace github.com/ethereum/go-ethereum => /home/user/go/src/github.com/ethereum/go-ethereum

//replace github.com/holiman/goevmlab => /home/user/go/src/github.com/holiman/goevmlab
