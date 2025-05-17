### Corpus conversion

This is a utility to convert between different corpus-formats;

- Golang-native corpus-format
- Libfuzzer corpus format
- Hexadecimal-encoded strings

## Golang-native corpus

The golang-native corpus is a set of files, one file per entry. Within the file,
the data is in the form of a golang instantiation, e.g:

```
go test fuzz v1
[]byte("\x03\xf8x00")
```
The files are typically named as the first 8 bytes of the `sha256` - hash of the file content (i.e the entire content).

## Libfuzzer corpus format

Libfuzzer also uses a set of files, one file per entry, but each file contains only the
raw binary data.

The files are typically named as the `SHA1`-hash of the vector.

## Hexadecimal-encoded strings

This does not use `in.path`, nor `out.path`, but reads from `stdin` and writes
to `stdout`.

# Examples

Converting from and to `golang`:
```
./corpusconvert -in.type "golang" -in.path ./sample.golang/ -out.type "golang" -out.path /tmp/
```
Same with `libfuzzer`:
```
./corpusconvert -in.type "libfuzzer" -in.path ./sample.libfuzzer/ -out.type "libfuzzer" -out.path /tmp/
```
And for `hex`
```
$ cat sample.hex | ./corpusconvert -in.type="hex" -out.type="hex"
0x1122
0x3344
0xffffaabbccdd
0x0badda
0x1337
```

