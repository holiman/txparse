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
	//	in := os.Args[1]
	// comma separated binaries
	binaries := os.Args[1]
	bins := strings.Split(binaries, ",")
	if len(binaries) < 2 {
		fmt.Printf("Usage: comparer parser1,parser2,... \n")
		fmt.Printf("Pipe input to process")
		return
	}
	if err := doit(bins); err != nil {
		fmt.Printf("err: %v", err)
	}
}

type proc struct {
	cmd    string
	outp   io.ReadCloser
	inp    io.WriteCloser
	outbuf *bufio.Scanner
}

func doit(bins []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	var procs []proc

	for _, bin := range bins {
		cmdArgs := strings.Split(bin, " ")
		var args []string
		if len(cmdArgs) > 1 {
			args = cmdArgs[1:]
		}
		cmd := exec.Command(cmdArgs[0], args...)
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
			cmd:    cmd.String(),
			outp:   stdout,
			inp:    stdin,
			outbuf: bufio.NewScanner(stdout),
		})
	}
	fmt.Printf("Using %d processes\n", len(procs))
	if len(procs) < 2 {
		return errors.New("At least 2 processes are needed")
	}
	for i, proc := range procs {
		fmt.Printf("  %d: %v\n", i, proc.cmd)
	}
	var count = 0
	var lastLog = time.Now()

	for scanner.Scan() {
		if time.Since(lastLog) > 8*time.Second {
			fmt.Fprintf(os.Stdout, "# %d cases OK\n", count)
			lastLog = time.Now()
		}
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
				fmt.Printf("process read failure: %v\n", proc.cmd)
				fmt.Printf("input: %v\n", l)
				return fmt.Errorf("process read fail: %v", proc.cmd)
			}
			outputs = append(outputs, cur)
			if i == 0 {
				prev = cur
				continue
			}
			if strings.HasPrefix(cur, "Exception") {
				// work-around nethermind thing
				cur = fmt.Sprintf("err: %v", cur)
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
			for j, outp := range outputs {
				fmt.Printf("%d: proc %d: %v\n", count, j, outp)
			}
			fmt.Printf("%d input %v\n", count, l)
			fmt.Fprintln(os.Stderr, l)
			fmt.Printf("\n")
		}
	}
	fmt.Fprintf(os.Stdout, "# %d cases OK", count)
	return nil
}
