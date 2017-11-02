package main

import (
	"debug/elf"
	"debug/gosym"
)

func getSymbolTable(prog string) *gosym.Table {
	exe, err := elf.Open(prog)
	if err != nil {
		panic(err)
	}
	defer exe.Close()

	addr := exe.Section(".text").Addr

	lineTableData, err := exe.Section(".gopclntab").Data()
	if err != nil {
		panic(err)
	}
	lineTable := gosym.NewLineTable(lineTableData, addr)
	if err != nil {
		panic(err)
	}

	symTableData, err := exe.Section(".gosymtab").Data()
	if err != nil {
		panic(err)
	}

	symTable, err := gosym.NewTable(symTableData, lineTable)
	if err != nil {
		panic(err)
	}

	return symTable
}
