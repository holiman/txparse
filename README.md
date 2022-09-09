### Tx parser

This is a very simple utility, which reads line by line from standard input. 
For each line, it tries to interpret it as hexadecimal data, and the data as
an Ethereum transaction (specifically: London-enabled and `chainid=1`).

If all goes well, it outputs a line containing the `address` of the sender. 
Otherwise, it outputs `err: ` and a suitable error message. 

Example: 

```
$ cat sample.input | ./txparse 
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
