## Custom fuzzers

This repository contains a few different (two at the moment) one-off bespoke fuzzing targets, for Ethereum 'stuff'. 

### Tx parser (`cmd/txparse`)

This is a very simple utility, which reads line by line from standard input.

1.Read a line of input
2. If the line starts with `#`, it is ignored, and no output is emitted. Go to 1.
3. Otherwise, remove any `non-alnum` characters from the input line (`[^0-9A-Za-z]`) 
4. Try to interpret it as hexadecimal data
4a. If error: yield error and goto 1.
5. Parse a transaction (cancun-enabled, chainid `1`) from the data.
5b. If error: yield error and goto 1.
6. Emit the transaction `sender`
7. Go to 1

If all goes well, it outputs a line containing the `address` of the sender. 
Otherwise, it outputs `err: ` and a suitable error message.

(First build with `cd txparse && go build .`)
Example:


```
$ cat sample.input | ./txparse 
err: typed transaction too short
err: typed transaction too short
err: typed transaction too short
err: rlp: expected input list for types.AccessListTx
err: rlp: value size exceeds available input length
err: rlp: input string too long for uint64, decoding into (types.LegacyTx).Nonce
0xd02d72e067e77158444ef2020ff2d325f929b363
0xd02d72e067e77158444ef2020ff2d325f929b363
0xd02d72e067e77158444ef2020ff2d325f929b363
err: transaction type not supported
```
- `./cmd/txparse/random_corpus.txt` contains a large batch of tests, combinations of fuzzing-corpi from various runs
- `./cmd/txparse/reth_corpus.txt` contains corpus from fuzzing on reth

## `corpusconvert`

`corpusconvert` is a small converter between hexadecimal text and fuzzing-corpus, the format
used by libfuzzer (golang native fuzzing uses a similar but different format! Also individual files, but the content are not just raw bytes, but the golang instantiation of the input). Each vector is an individual file containing
the raw bytes.


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
