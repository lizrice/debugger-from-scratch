package main

import (
	"debug/gosym"
	"encoding/binary"
	"fmt"
	"syscall"
)

func outputStack(symTable *gosym.Table, pid int, regs *syscall.PtraceRegs) {

	// fmt.Printf("%#v\n", regs)

	_, _, fn = symTable.PCToLine(regs.Rip)

	sp := regs.Rsp
	bp := regs.Rbp
	var i uint64
	var nextbp uint64

	for {
		i = 0
		frameSize := bp - sp + 8

		if frameSize > 1000 || bp == 0 {
			fmt.Printf("Strange frame size: SP: %X | BP : %X \n", sp, bp)
			return
			// frameSize = 32
			// bp = sp + frameSize - 8
		}

		b := make([]byte, frameSize)
		_, err := syscall.PtracePeekData(pid, uintptr(sp), b)
		if err != nil {
			panic(err)
		}

		content := binary.LittleEndian.Uint64(b[i : i+8])
		_, lineno, nextfn := symTable.PCToLine(content)
		if nextfn != nil {
			fn = nextfn
			// fmt.Printf("  %X %X: return to %s line %d\n", sp, content, fn.Name, lineno)
			fmt.Printf("  called by %s line %d\n", fn.Name, lineno)
		}

		for i = 8; sp+i <= bp; i += 8 {
			content := binary.LittleEndian.Uint64(b[i : i+8])
			if sp+i == bp {
				// fmt.Printf("Frame added calling %s\n", fn.Name)
				nextbp = content
			}
			// fmt.Printf("  %X %X  \n", sp+i, content)
		}

		if fn.Name == "main.main" || fn.Name == "runtime.main" {
			break
		}

		sp = sp + i
		bp = nextbp
	}

	fmt.Println()
}
