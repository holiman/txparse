### Comparer

This is a utility to compare the output of running several `txparse`-like 
binaries against each other. It basically multiplexes the input lines to 
a number of binaries, and compares the outputs. 

Example: 
```
cat ../txparse/random_corpus.txt | go run . ./binaries.txt
Processes:
0: /home/user/workspace/txparse/txparse
1: /home/user/workspace/besu/build/install/besu/bin/besu txparse

# 498 cases OK
```

Where

- `../txparse/random_corpus.txt` contains `498` lines of input
- `./binaries.txt` contains info on how to start each process. OBS: This may be a
  composite command, such as `docker run ...`. See example below

```
# geth
/home/user/go/src/github.com/holiman/txparse/txparse
# besu
/home/user/workspace/besu/build/install/besu/bin/besu txparse
```