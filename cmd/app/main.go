package main

import (
	"lc3-vm/internal/vm"
)

func main() {
	v := vm.NewVM()

	if err := v.LoadRom("2048.obj"); err != nil {
		panic(err)
	}

	v.Run()
}
