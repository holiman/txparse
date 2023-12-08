### Tx parser (`cmd/txparse`)

This is a very simple utility, which reads line by line from standard input.

1. Read a line of input
2. If the line starts with `#`, it is ignored, and no output is emitted. Go to 1.
3. Otherwise, remove any `non-alnum` characters from the input line (`[^0-9A-Za-z]`)
4. Try to interpret it as hexadecimal data
   4a. If error: yield error and goto 1.
5. Parse a transaction (cancun-enabled, chainid `1`, consensus-encoding (== no blob sidecars)) from the data.
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

### Corpus

The files `random_corpus.txt` is updated from time to time, to contain everything that increases the code coverage.

### Corpus input 

To use golang native fuzzing to increase the corpus, first run it: 

```
$ GOCACHE=`pwd`/corpus  go test . -fuzz Fuzz
```

This will create corpus into the folder `./corpus`. When you are finished, it's time to extract the 
corpus again, from the golang native format into hex. For this, we use `cmd/corpustoinput`

```
go run ../corpustoinput ./corpus/fuzz/github.com/holiman/txparse/cmd/txparse/Fuzz/ > new_corpus.txt
```
And finally, merge with the old
```
cat new_corpus.txt random_corpus.txt | sort | uniq > tmpfile
mv  tmpfile random_corpus.txt 
rm  new_corpus.txt
```

#### Output

THe file `random_corpus.output` is regenerated along with the input: it can be used as a base to sanity-check
a txparse-implementation. It's generated like this: 

```
cat random_corpus.txt | go run . | tee random_corpus.output
```
