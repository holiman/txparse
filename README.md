## Custom fuzzers

This repository contains a few different (two at the moment) one-off bespoke fuzzing targets, for Ethereum 'stuff'. 

### Tx parser (`txparse`)

This is a very simple utility, which reads line by line from standard input.
For each line, it tries to interpret it as hexadecimal data, and the data as
an Ethereum transaction (specifically: London-enabled and `chainid=1`).

If all goes well, it outputs a line containing the `address` of the sender.
Otherwise, it outputs `err: ` and a suitable error message.

(First build with `cd txparse && go build .`)
Example:


```
$ cat ./sample.input | ./txparse 
err: typed transaction too short
err: typed transaction too short
err: typed transaction too short
err: transaction type not supported
err: rlp: value size exceeds available input length
err: rlp: input string too long for uint64, decoding into (types.LegacyTx).Nonce
0xd02d72e067e77158444ef2020ff2d325f929b363
0xd02d72e067e77158444ef2020ff2d325f929b363
err: transaction type not supported
```

### EOF parser (`eofparse`)

This is a very simple utility, which reads line by line from standard input.
For each line, it tries to interpret it as hexadecimal data, and the data as
an EOF1 code blob (verified using the Shanghai jumptable).

If all goes well, it outputs a line containing `OK ` followed by the comma-separted hex-encoded code sections.

Example: 
```
OK 604200,6042604200,00
```
Otherwise, it outputs `err: ` and a suitable error message.
Example: 
```
err: use of undefined opcode opcode 0x22 not defined
```

It requires the use of go-ethereum on the `lightclient/eof` branch:

Example:


```
]$  cat ./all.input | ./eofparse | head -n 10
err: invalid magic
err: invalid magic
err: container size less than minimum valid size
err: container size less than minimum valid size
err: container size less than minimum valid size
err: container size less than minimum valid size
err: container size less than minimum valid size
err: container size less than minimum valid size
err: container size less than minimum valid size
err: container size less than minimum valid size
```

A larger corpus is in `all.input`. That also has the corresponding output, generated via 
`cat all.input | ./eofparse > all.output`. So you can use `all.output` to compare. But remember, 
anything after `err:` is up to the implementation to phrase how they see fit. 
