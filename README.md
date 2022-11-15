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

If all goes well, it outputs a line containing `OK ` followed by the length of the EOF header.
Otherwise, it outputs `err: ` and a suitable error message.

(First build with `cd eof && go build .`)
Example:


```
]$ cat ./sample.input | ./eofparse 
err: invalid version byte
err: invalid version byte
OK 7
OK 7
err: unknown section id
err: invalid version byte
err: invalid version byte
OK 7
OK 7
OK 10
OK 10
OK 10
OK 10
err: code section doesn't end with terminating instruction: PUSH1
OK 7
OK 7
OK 7
OK 10
err: invalid version byte
err: no code section
err: no code section
err: can't read code section size
err: can't read code section size
err: invalid total size
err: invalid total size
err: invalid total size
err: code section size is 0
err: code section size is 0
err: data section before code section
err: data section before code section
err: can't read data section size
err: can't read data section size
err: invalid total size
err: invalid total size
err: invalid total size
err: invalid total size
err: data section size is 0
err: multiple data sections
err: unknown section id
err: undefined instruction: opcode 0xc not defined
err: undefined instruction: opcode 0xef not defined
err: code section doesn't end with terminating instruction: ADDRESS
err: code section doesn't end with terminating instruction: PUSH1
err: code section doesn't end with terminating instruction: PUSH32
err: code section doesn't end with terminating instruction: PUSH32
err: relative offset points to immediate argument
err: relative offset points to immediate argument
err: relative offset points to immediate argument
err: undefined instruction: opcode 0xc not defined
err: undefined instruction: opcode 0xc not defined
err: undefined instruction: opcode 0xc not defined
```

A larger corpus is in `all.input`. That also has the corresponding output, generated via 
` cat all.input | ./eofparse > all.output`. So you can use `all.output` to compare. But remember, 
anything after `err:` is up to the implementation to phrase how they see fit. 
