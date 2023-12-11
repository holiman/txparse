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

func setup() ([]string, error) {
	var binaries []string
	if len(os.Args) < 2 {
		fmt.Printf("Usage: comparer <file with binaries> \n")
		fmt.Printf("Pipe input to process\n")
		return nil, errors.New("insufficient arguments")
	}
	binFile, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}
	defer binFile.Close()
	scanner := bufio.NewScanner(binFile)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		if len(strings.TrimSpace(scanner.Text())) == 0 {
			continue
		}
		binaries = append(binaries, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(binaries) < 2 {
		fmt.Printf("You need to provide at least two binaries to use (have %d)\n", len(binaries))
		return nil, errors.New("insufficient binaries")
	}
	return binaries, nil
}

func main() {
	bins, err := setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	var input = make(chan string)
	go func() {
		defer close(input)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, 200_000), 1_000_000)
		for scanner.Scan() {
			input <- scanner.Text()
		}
	}()
	if err := doit(bins, input, nil); err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
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
	outbuf := bufio.NewScanner(stdout)
	outbuf.Buffer(make([]byte, 200_000), 1_000_000)
	return &proc{
		cmd:    cmd.String(),
		outp:   stdout,
		inp:    stdin,
		outbuf: outbuf,
	}, nil

}
func doit(bins []string, inputs chan string, results chan string) error {
	var procs []*proc
	// Start up the processes
	for _, bin := range bins {
		p, err := startProcess(strings.Split(bin, " "))
		if err != nil {
			return fmt.Errorf("could not start process %q: %v", bin, err)
		}
		procs = append(procs, p)
	}
	if len(procs) < 2 {
		return errors.New("At least 2 processes are needed")
	}
	fmt.Printf("Processes:\n")
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
			fmt.Fprintf(errMsg, "Processes:\n")
			for i, proc := range procs {
				fmt.Fprintf(errMsg, "  %d: %v\n", i, proc.cmd)
			}
			fmt.Fprintf(errMsg, "\n")
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
	fmt.Fprintf(os.Stdout, "# %d cases OK\n", count)
	return nil
}
