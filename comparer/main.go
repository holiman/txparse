package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	in := os.Args[1]
	// comma separated binaries
	binaries := os.Args[2]
	bins := strings.Split(binaries, ",")
	if err := doit(bins, in); err != nil {
		fmt.Printf("err: %v", err)
	}
}

type proc struct {
	outp   io.ReadCloser
	inp    io.WriteCloser
	outbuf *bufio.Scanner
}

func doit(bins []string, input string) error {
	var procs []proc
	indata, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	var (
		scanner = bufio.NewScanner(indata)
	)

	for _, bin := range bins {
		cmd := exec.Command(bin)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		if err = cmd.Start(); err != nil {
			return err
		}
		procs = append(procs, proc{
			outp:   stdout,
			inp:    stdin,
			outbuf: bufio.NewScanner(stdout),
		})
	}
	var count = 0
	for scanner.Scan() {
		count++
		l := scanner.Text()
		for _, proc := range procs {
			proc.inp.Write([]byte(l))
			proc.inp.Write([]byte("\n"))
		}
		prev := ""
		var ok = true
		var outputs []string
		for i, proc := range procs {
			var cur = ""
			if proc.outbuf.Scan() {
				cur = proc.outbuf.Text()
			} else {
				panic("foo")
			}
			outputs = append(outputs, cur)
			if i == 0 {
				prev = cur
				continue
			}
			if strings.HasPrefix(prev, "err") && strings.HasPrefix(cur, "err") {
				prev = cur
				continue
			}
			if prev != cur {
				ok = false
			}
			prev = cur
		}
		if !ok || true {
			for j, outp := range outputs {
				fmt.Printf("%d: proc %d: %v\n", count, j, outp)
			}
			fmt.Printf("%d input %v\n", count, l)
			fmt.Fprintln(os.Stderr, l)
		}
	}
	return nil
}
