package vm

import (
	"fmt"
	"syscall"
)

const (
	MR_KBSR = 0xFE00 /* keyboard status */
	MR_KBDR = 0xFE02 /* keyboard data */
)

func checkKey() uint16 {
	var readfds syscall.FdSet

	readfds.Bits[0] = 1 << (uint(syscall.Stdin) % (8 * syscall.FD_SETSIZE))
	if err := syscall.Select(int(syscall.Stdin+1), &readfds, nil, nil, nil); err != nil {
		return 0
	}

	return uint16(len(readfds.Bits))
}

func getChar() uint16 {
	var input byte
	_, err := fmt.Scanf("%c", &input)
	if err != nil {
		panic(err)
	}

	return uint16(input & 0xFF)
}

func (vm *Vm) memWrite(address, value uint16) {
	vm.register.general[address] = value
}

func (vm *Vm) memRead(address uint16) uint16 {
	if address == MR_KBSR {
		if checkKey() == 1 {
			vm.memory[MR_KBSR] = (1 << 15)
			vm.memory[MR_KBDR] = getChar()
		} else {
			vm.memory[MR_KBSR] = 0
		}
	}

	return vm.memory[address]
}
