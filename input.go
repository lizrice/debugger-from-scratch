package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func inputContinue(pid int) bool {
	sub := false
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n(C)ontinue, (S)tep, set (B)reakpoint or (Q)uit? > ")
	for {
		scanner.Scan()
		input := scanner.Text()
		switch strings.ToUpper(input) {
		case "C":
			return true
		case "S":
			return false
		case "B":
			fmt.Printf("  Enter line number in %s: > ", targetfile)
			sub = true
		case "Q":
			os.Exit(0)
		default:
			if sub {
				line, _ = strconv.Atoi(input)
				breakpointSet, originalCode = setBreak(pid, targetfile, line)
				return true
			}
			fmt.Printf("Unexpected input %s\n", input)
			fmt.Printf("\n(C)ontinue, (S)tep, set (B)reakpoint or (Q)uit? > ")
		}
	}
}

func setBreak(pid int, filename string, line int) (bool, []byte) {
	var err error
	pc, _, err = symTable.LineToPC(filename, line)
	if err != nil {
		fmt.Printf("Can't find breakpoint for %s, %d\n", filename, line)
		return false, []byte{}
	}

	// fmt.Printf("Stopping at %X\n", pc)
	return true, replaceCode(pid, pc, []byte{0xCC})
}
