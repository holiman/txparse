package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	binaries := os.Args[1]
	bins := strings.Split(binaries, ",")
	if len(binaries) < 2 {
		fmt.Printf("Usage: comparer parser1,parser2,... \n")
		fmt.Printf("Pipe input to process")
		return
	}
	var input = make(chan string)
	go func() {
		defer close(input)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input <- scanner.Text()
		}
	}()
	if err := doit(bins, input, nil); err != nil {
		fmt.Printf("err: %v", err)
	}
}

type proc struct {
	cmd    string
	outp   io.ReadCloser
	inp    io.WriteCloser
	outbuf *bufio.Scanner
}

func startProcess(cmdArgs []string) (*proc, error) {
	var args []string
	if len(cmdArgs) > 1 {
		args = cmdArgs[1:]
	}
	for _, arg := range args {
		if len(arg) == 0 {
			// probably a double-space
			fmt.Printf("Warn: empty arg (double-space?) in '%v'\n", cmdArgs)
		}
	}
	cmd := exec.Command(cmdArgs[0], args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}
	return &proc{
		cmd:    cmd.String(),
		outp:   stdout,
		inp:    stdin,
		outbuf: bufio.NewScanner(stdout),
	}, nil

}
func doit(bins []string, inputs chan string, results chan string) error {
	var procs []*proc
	// Start up the processes
	for _, bin := range bins {
		p, err := startProcess(strings.Split(bin, " "))
		if err != nil {
			return err
		}
		procs = append(procs, p)
	}
	if len(procs) < 2 {
		return errors.New("At least 2 processes are needed")
	}
	fmt.Printf("Using %d processes:\n", len(procs))
	for i, proc := range procs {
		fmt.Printf("  %d: %v\n", i, proc.cmd)
	}
	fmt.Println("")
	var (
		count   = 0
		lastLog = time.Now()
	)
	for l := range inputs {
		if len(l) == 0 || strings.HasPrefix(l, "#") {
			if results != nil {
				results <- ""
			}
			continue
		}
		if time.Since(lastLog) > 8*time.Second {
			fmt.Fprintf(os.Stdout, "# %d cases OK\n", count)
			lastLog = time.Now()
		}
		count++
		// Feed inputs
		for _, proc := range procs {
			proc.inp.Write([]byte(l))
			proc.inp.Write([]byte("\n"))
		}
		var (
			prev    = ""
			ok      = true
			outputs []string
		)
		// Compare outputs
		for i, proc := range procs {
			var cur = ""
			if proc.outbuf.Scan() {
				cur = proc.outbuf.Text()
			} else {
				err := proc.outbuf.Err()
				a := fmt.Sprintf("%d: process read failure: %v %v\ninput: %v\n", count, proc.cmd, err, l)
				fmt.Printf(a)
				if results != nil {
					results <- a
				}
				return fmt.Errorf("process read fail line %d: %v", count, proc.cmd)
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
		if !ok {
			var errMsg = new(strings.Builder)
			fmt.Fprintf(errMsg, "%d input %v\n", count, l)
			for j, outp := range outputs {
				fmt.Fprintf(errMsg, "%d: proc %d: %v\n", count, j, outp)
			}
			fmt.Fprintf(errMsg, "\n")
			fmt.Printf(errMsg.String())
			fmt.Fprintln(os.Stderr, l)
			if results != nil {
				results <- errMsg.String()
			}
		} else {
			if results != nil {
				results <- ""
			}
		}
	}
	fmt.Fprintf(os.Stdout, "# %d cases OK", count)
	return nil
}
