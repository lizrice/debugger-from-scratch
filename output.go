package main

import (
	"debug/gosym"
	"encoding/binary"
	"fmt"
	"syscall"
)

func outputStack(symTable *gosym.Table, pid int, ip uint64, sp uint64, bp uint64) {

	_, _, fn = symTable.PCToLine(ip)

	var i uint64
	var nextbp uint64

	for {
		i = 0
		frameSize := bp - sp + 8

		// If we look at bp / sp while they are being updated we can
		// get some odd results
		if frameSize > 1000 || bp == 0 {
			fmt.Printf("Strange frame size: SP: %X | BP : %X \n", sp, bp)
			frameSize = 32
			bp = sp + frameSize - 8
		}

		// Read the next stack frame
		b := make([]byte, frameSize)
		_, err := syscall.PtracePeekData(pid, uintptr(sp), b)
		if err != nil {
			panic(err)
		}

		// The address to return to is at the top of the frame
		content := binary.LittleEndian.Uint64(b[i : i+8])
		_, lineno, nextfn := symTable.PCToLine(content)
		if nextfn != nil {
			fn = nextfn
			fmt.Printf("  called by %s line %d\n", fn.Name, lineno)
		}

		// Rest of the frame
		for i = 8; sp+i <= bp; i += 8 {
			content := binary.LittleEndian.Uint64(b[i : i+8])
			if sp+i == bp {
				nextbp = content
			}
			// fmt.Printf("  %X %X  \n", sp+i, content)
		}

		// We want to stop the stack trace at main.main.
		// At the point where bp & sp are being updated we may
		// miss main.main, so we backstop with runtime.main.
		if fn.Name == "main.main" || fn.Name == "runtime.main" {
			break
		}

		// Move to the next frame
		sp = sp + i
		bp = nextbp
	}

	fmt.Println()
}
