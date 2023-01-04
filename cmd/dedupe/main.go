package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	work()
}

func work() {
	var (
		vectors  = make(map[string]bool)
		scanner  = bufio.NewScanner(os.Stdin)
		toRemove = regexp.MustCompile(`[^0-9A-Za-z]`)
	)
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			fmt.Println(l)
			continue
		}
		sanitized := toRemove.ReplaceAllString(l, "")
		sanitized = strings.ToLower(sanitized)
		if vectors[sanitized] { // already present
			//fmt.Printf("#dup: %v\n", l)
			continue
		}
		vectors[sanitized] = true
		fmt.Println(l)
	}
}
