package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

var (
	input   = flag.String("in.type", "", "Input-type: 'hex', 'libfuzzer' or 'golang'")
	output  = flag.String("out.type", "", "Output-type: 'hex', 'libfuzzer' or 'golang'")
	inpath  = flag.String("in.path", ".", "path to read corpus-files from")
	outpath = flag.String("out.path", ".", "path to write corpus-files from")
)

func main() {
	flag.Parse()
	if err := doit(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// readerFn is a function which can read and emit inputs
type readerFn func() ([]byte, error)

// writerFn is a function which writes data to an output
type writerFn func([]byte) error

func doit() error {
	reader, err := configureReader()
	if err != nil {
		return err
	}
	writer, err := configureWriter()
	if err != nil {
		return err
	}
	for i := 0; ; i++ {
		data, err := reader()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error reading vector %d: %v\n", i, err)
		}
		if err = writer(data); err != nil {
			return err
		}
	}
	return nil
}

func configureReader() (readerFn, error) {
	var reader readerFn
	switch *input {
	case "hex":
		var (
			scanner  = bufio.NewScanner(os.Stdin)
			toRemove = regexp.MustCompile(`[^0-9A-Za-z]`)
		)
		scanner.Buffer(make([]byte, 1024*1024), 5*1024*1024)
		reader = func() ([]byte, error) {
			for scanner.Scan() {
				if bytes.HasPrefix(scanner.Bytes(), []byte("#")) {
					continue
				}
				sanitized := toRemove.ReplaceAll(scanner.Bytes(), []byte{})
				return common.FromHex(string(sanitized)), nil
			}
			return nil, io.EOF
		}
	case "libfuzzer":
		var (
			i            = 0
			entries, err = os.ReadDir(*inpath)
		)
		if err != nil {
			return nil, err
		}
		reader = func() ([]byte, error) {
			for i < len(entries) {
				entry := entries[i]
				i++
				if entry.IsDir() {
					continue
				}
				path := filepath.Join(*inpath, entry.Name())
				if data, err := os.ReadFile(path); err != nil {
					return nil, fmt.Errorf("failed to read file %v: %v", path, err)
				} else {
					return data, nil
				}
			}
			return nil, io.EOF
		}
	case "golang":
		var (
			i            = 0
			entries, err = os.ReadDir(*inpath)
		)
		if err != nil {
			return nil, err
		}
		reader = func() ([]byte, error) {
			for i < len(entries) {
				entry := entries[i]
				i++
				if entry.IsDir() {
					continue
				}
				path := filepath.Join(*inpath, entry.Name())
				data, err := os.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("failed to read file %v: %v", path, err)
				}
				lines := bytes.Split(data, []byte("\n"))
				if len(lines) < 2 {
					return nil, fmt.Errorf("must include version and at least one value")
				}
				if len(lines) > 3 {
					return nil, fmt.Errorf("why so many? (file %v)", path)
				}
				for _, line := range lines[1:] {
					line = bytes.TrimSpace(line)
					if len(line) == 0 {
						continue
					}
					v, err := parseCorpusValue(line)
					if err != nil {
						return nil, fmt.Errorf("malformed line %q: %v", line, err)
					}
					if data, ok := v.([]byte); ok {
						return data, nil
					} else {
						return nil, fmt.Errorf("unsupported fuzzing-type %T", v)
					}
				}
				return nil, fmt.Errorf("unexpected eof in %v", path)
			}
			return nil, io.EOF
		}
	default:
		return nil, fmt.Errorf("-input must be one of 'hex', 'libfuzzer' or 'golang'")
	}
	return reader, nil
}

func configureWriter() (writerFn, error) {
	var writer writerFn

	switch *output {
	case "hex":
		writer = func(data []byte) error {
			fmt.Printf("%#x\n", data)
			return nil
		}
	case "libfuzzer":
		finfo, err := os.Stat(*outpath)
		if err != nil {
			return nil, err
		}
		if !finfo.IsDir() {
			return nil, fmt.Errorf("not a directory: %v", *outpath)
		}
		writer = func(data []byte) error {
			fname := fmt.Sprintf("%x", sha1.Sum(data))
			path := filepath.Join(*outpath, fname)
			if err := os.WriteFile(path, data, 0777); err != nil {
				return err
			}
			return nil
		}
	case "golang":
		finfo, err := os.Stat(*outpath)
		if err != nil {
			return nil, err
		}
		if !finfo.IsDir() {
			return nil, fmt.Errorf("not a directory: %v", *outpath)
		}
		writer = func(data []byte) error {
			formatted := fmt.Sprintf("go test fuzz v1\n[]byte(%q)\n", data)
			sum := fmt.Sprintf("%x", sha256.Sum256([]byte(formatted)))[:16]
			path := filepath.Join(*outpath, sum)
			if err := os.WriteFile(path, []byte(formatted), 0777); err != nil {
				return err
			}
			return nil
		}
	default:
		return nil, fmt.Errorf("-output must be one of 'hex', 'libfuzzer' or 'golang'")
	}
	return writer, nil
}

// parseCorpusValue is taken mostly from the go internals, in internal/fuzz/encoding.go.
// It has been stripped out, so it only supports []byte types
func parseCorpusValue(line []byte) (any, error) {
	fs := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fs, "(test)", line, 0)
	if err != nil {
		return nil, err
	}
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil, fmt.Errorf("expected call expression")
	}
	if len(call.Args) != 1 {
		return nil, fmt.Errorf("expected call expression with 1 argument; got %d", len(call.Args))
	}
	arg := call.Args[0]

	if arrayType, ok := call.Fun.(*ast.ArrayType); ok {
		if arrayType.Len != nil {
			return nil, fmt.Errorf("expected []byte or primitive type")
		}
		elt, ok := arrayType.Elt.(*ast.Ident)
		if !ok || elt.Name != "byte" {
			return nil, fmt.Errorf("expected []byte")
		}
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return nil, fmt.Errorf("string literal required for type []byte")
		}
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	}
	return nil, fmt.Errorf("Unsupported type")
}
