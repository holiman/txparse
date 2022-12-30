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

If all goes well, it outputs a line containing `OK ` followed by the hex-encoded bytecode of the first code section. 
Otherwise, it outputs `err: ` and a suitable error message.

It requires a change to go-ethereum, which also needs to run on `lightclient/eof` branch:
```diff
diff --git a/core/vm/jump_table.go b/core/vm/jump_table.go
index a91406ef6e..7db9b1b95c 100644
--- a/core/vm/jump_table.go
+++ b/core/vm/jump_table.go
@@ -85,6 +85,10 @@ func validate(jt JumpTable) JumpTable {
        return jt
 }
 
+func NewShanghaiEOFInstructionSetForTesting() JumpTable {
+       return newShanghaiEOFInstructionSet()
+}
+
 func newShanghaiEOFInstructionSet() JumpTable {
        instructionSet := newMergeInstructionSet()
        enable3855(&instructionSet)
```


(First build with `cd eof && go build .`)
Example:


```
$ cat ./sample.input | ./eofparse | head -n 10
err: invalid version byte
err: invalid version byte
OK 00
OK 00
err: unknown section id
err: invalid version byte
err: invalid version byte
OK 00
OK 0000
OK 0000
```

A larger corpus is in `all.input`. That also has the corresponding output, generated via 
`cat all.input | ./eofparse > all.output`. So you can use `all.output` to compare. But remember, 
anything after `err:` is up to the implementation to phrase how they see fit. 
